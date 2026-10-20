package dashboard

const personalJS = `(function(){
'use strict';
const key='newsgo.subscriptions.v2';
const legacyPersonalKey='newsgo.personalFeeds.v1';
const legacyIncludedKey='newsgo.includedSubscriptions.v1';
const catalog=Array.from(document.querySelectorAll('[data-included-source]'),b=>({name:b.dataset.includedSource,url:b.dataset.feedUrl,category:b.dataset.category||'Other'}));
const catalogByName=new Map(catalog.map(f=>[f.name,f]));
const feedIdentity=u=>{try{const x=new URL(u);x.hash='';return 'feed:'+x.href}catch(_){return 'feed:'+u}};
const sharedRecord=name=>({id:'shared:'+name,kind:'shared',name});
const feedRecord=f=>({id:feedIdentity(f.url),kind:'feed',name:f.name,url:f.url,category:f.category||'Other'});
const recordValid=r=>r&&typeof r.id==='string'&&typeof r.name==='string'&&r.name.length>0&&r.name.length<=140&&
 (r.kind==='shared'||(r.kind==='feed'&&typeof r.url==='string'&&/^https?:\/\/[^\s/]+/i.test(r.url)&&r.url.length<=2048));
// Publisher exclusions are independent of subscription membership.
const exclusionKey='newsgo.excludedDomains.v1';
let excludedDomains=[],exclusionError='';
function hostOf(value){
 try{const u=new URL(value.includes('://')?value:'https://'+value);if(!/^https?:$/.test(u.protocol)||u.username||u.password)return '';return u.hostname.toLowerCase().replace(/\.$/,'').replace(/^www\./,'')}
 catch(_){return ''}
}
function exclusionDomain(value){
 const host=hostOf(value);
 return host&&/^[a-z0-9.-]+$/.test(host)&&host.includes('.')&&!host.startsWith('.')&&!host.endsWith('.')&&!host.includes('..')?host:'';
}
function excludedURL(value){const h=hostOf(value);return !!h&&excludedDomains.some(d=>h===d||h.endsWith('.'+d))}
function readExclusions(){try{const data=JSON.parse(localStorage.getItem(exclusionKey)||'[]');if(!Array.isArray(data)||data.some(x=>typeof x!=='string'||exclusionDomain(x)!==x))throw Error('Invalid excluded-domain list');excludedDomains=[...new Set(data)]}catch(e){exclusionError='Exclusions were not changed: '+e.message}}
function saveExclusions(list){
 try{const unique=[...new Set(list)];localStorage.setItem(exclusionKey,JSON.stringify(unique));excludedDomains=unique;refreshExclusionViews();return true}
 catch(e){message('Could not save exclusions: '+e.message);return false}
}
function excludeDomain(value){const domain=exclusionDomain(value);if(!domain){message('Enter a valid website domain.');return}if(!excludedDomains.includes(domain))saveExclusions([...excludedDomains,domain]);else message(domain+' is already excluded.')}
function toggleDomain(value){const domain=exclusionDomain(value);if(!domain){message('Could not identify a website domain.');return}if(excludedDomains.includes(domain))saveExclusions(excludedDomains.filter(x=>x!==domain));else excludeDomain(domain)}
function publisherForFeed(feed){return exclusionDomain(feed.website||feed.url||'')}
function statusLabel(domain){return domain&&excludedURL(domain)?' · Excluded':''}
function excludeButton(domain){const d=exclusionDomain(domain);if(!d)return null;const b=button(excludedURL(d)?'Allow this site':'Exclude this site',()=>toggleDomain(d),'exclude-site');b.dataset.publisherDomain=d;return b}
function markSite(info,domain){if(!domain)return;const badge=document.createElement('small');badge.className='exclusion-label';badge.textContent=statusLabel(domain);info.append(badge)}
function renderExclusionManager(){
 const root=el('excluded-list');root.replaceChildren();
 el('excluded-count').textContent='('+excludedDomains.length+')';
 if(!excludedDomains.length){root.textContent='No sites excluded.';return}
 excludedDomains.slice().sort().forEach(d=>{const row=document.createElement('div');row.className='personal-row';const label=document.createElement('strong');label.textContent=d;row.append(label,button('Allow site',()=>toggleDomain(d)));root.append(row)})
}
function selectedPublisherExcluded(){
 const value=el('read-source').value;
 if(value.startsWith('personal:')){const source=feeds.find(f=>f.url===value.slice(9));return !!source&&excludedURL(publisherForFeed(source))}
 return false;
}
function refreshExclusionViews(){
 // If the user just excluded the publication they are reading, show the
 // remaining subscriptions rather than an empty publisher page.
 if(selectedPublisherExcluded()){el('read-source').value='';selectedSource='';showView('news')}
 renderExclusionManager();renderManager();searchCatalog();
}
readExclusions();
let registry=[], storageError='';
const unifiedKey='newsgo.subscriptions.unified.v1';
function loadSubscriptions(){
 try{
  const current=localStorage.getItem(unifiedKey)||localStorage.getItem(key);
  if(current!==null){
   const value=JSON.parse(current);
   if(!Array.isArray(value))throw Error('Subscription registry has an invalid format.');
   if(value.some(r=>!recordValid(r)))throw Error('Subscription registry contains invalid entries; no data was discarded.');
   registry=convertLegacy(value);
   localStorage.setItem(unifiedKey,JSON.stringify(registry));
   return;
  }
  const oldFeeds=JSON.parse(localStorage.getItem(legacyPersonalKey)||'[]');
  const oldIncluded=JSON.parse(localStorage.getItem(legacyIncludedKey)||'[]');
  if(!Array.isArray(oldFeeds)||!Array.isArray(oldIncluded))throw Error('Legacy subscription data is invalid.');
  const proposed=[
   ...oldIncluded.filter(n=>typeof n==='string').map(sharedRecord),
   ...oldFeeds.map(feedRecord)
  ];
  if(proposed.some(r=>!recordValid(r)))throw Error('Legacy records require review; migration was not performed.');
  const unique=new Map(proposed.map(r=>[r.id,r]));
  const migrated=[...unique.values()];
  localStorage.setItem(key,JSON.stringify(migrated)); // Migration commits only after a successful write.
  registry=convertLegacy(migrated);
  localStorage.setItem(unifiedKey,JSON.stringify(registry));
 }catch(e){storageError='Subscription storage: '+e.message+' Existing records were left unchanged.'}
}
// Catalog publisher subscriptions migrate to identical feed records. Unknown
// legacy publishers are preserved for review instead of being discarded.
function convertLegacy(records){
 const converted=records.map(r=>{
  if(r.kind!=='shared')return r;
  const src=catalogByName.get(r.name);
  return src&&src.url?feedRecord(src):r;
 });
 return [...new Map(converted.map(r=>[r.id,r])).values()];
}
loadSubscriptions();
let included=[], feeds=[];
function syncDerived(){
 included=registry.filter(r=>r.kind==='shared').map(r=>r.name); // unresolved legacy records only
 feeds=registry.filter(r=>r.kind==='feed').map(r=>({name:r.name,url:r.url,category:r.category||'Other'}));
}
syncDerived();
function saveSubscriptions(next){
 try{
  const serialized=JSON.stringify(next);
  localStorage.setItem(unifiedKey,serialized);
  registry=next;syncDerived();return true;
 }catch(e){message('Could not save subscriptions. No changes were applied: '+e.message);return false}
}
function hasSubscription(id){return registry.some(r=>r.id===id)}
function subscribeRecord(record){if(hasSubscription(record.id))return true;return saveSubscriptions([...registry,record])}
function unsubscribeRecord(id){return saveSubscriptions(registry.filter(r=>r.id!==id))}
const suggestions=[
{name:'NASA News',category:'Science',url:'https://www.nasa.gov/news-release/feed/'},
{name:'NASA Science',category:'Science',url:'https://science.nasa.gov/feed/'},
{name:'NASA JPL',category:'Science',url:'https://www.jpl.nasa.gov/feeds/news/'},
{name:'NOAA News',category:'Science',url:'https://www.noaa.gov/rss.xml'},
{name:'Ars Technica Science',category:'Science',url:'https://feeds.arstechnica.com/arstechnica/science'}
];
const el=id=>document.getElementById(id);
const valid=f=>f&&typeof f.name==='string'&&f.name.trim().length>0&&f.name.length<=70&&typeof f.category==='string'&&f.category.length<=40&&typeof f.url==='string'&&f.url.length<=2048&&/^https?:\/\/[^\s/]+/i.test(f.url);
let selectedSource='';

function message(t){el('personal-message').textContent=t}
function button(label,handler,cls){const b=document.createElement('button');b.type='button';b.textContent=label;if(cls)b.className=cls;b.addEventListener('click',handler);return b}
let searchTimer, searchSeq=0;
async function searchCatalog(){
 const q=el('catalog-search').value.trim(),root=el('catalog-list'),status=el('catalog-message');
 const seq=++searchSeq;root.replaceChildren();
 if(q.length<2){status.textContent='Enter at least two characters to search for news sources.';return}
 status.textContent='Searching for “'+q+'”…';
 try{
  const res=await fetch('/api/search-feeds?q='+encodeURIComponent(q),{cache:'no-store'});
  const data=await res.json();if(seq!==searchSeq)return;
  if(!res.ok)throw Error(data.error||'Search unavailable (HTTP '+res.status+').');
  if(!data.feeds||!data.feeds.length){status.textContent='No matching sources found. Try another search or Add a website below.';return}
  data.feeds.forEach(f=>{
   const row=document.createElement('div');row.className='catalog-card';
   const info=document.createElement('div'),name=document.createElement('strong'),detail=document.createElement('small');
   name.textContent=f.name;detail.textContent=f.website||new URL(f.url).hostname;
   info.append(name,detail);const domain=publisherForFeed(f);markSite(info,domain);
   const followed=feeds.some(x=>x.url===f.url);
   const b=button(followed?'Subscribed':'Subscribe',()=>add({name:f.name,category:'Other',url:f.url},b));b.dataset.feedUrl=f.url;b.disabled=followed;
   const exclude=excludeButton(domain);row.append(info,b);if(exclude)row.append(exclude);root.append(row);
  });
  status.textContent=data.feeds.length+' candidate sources. Each feed is validated when you select Subscribe.';
 }catch(err){if(seq===searchSeq)status.textContent='Search failed: '+err.message+' Try Add a website below.'}
}
function queueSearch(){clearTimeout(searchTimer);searchTimer=setTimeout(searchCatalog,350)}
function syncFollowButtons(){document.querySelectorAll('[data-feed-url]').forEach(b=>{const following=feeds.some(f=>f.url===b.dataset.feedUrl);b.textContent=following?'Subscribed':'Subscribe';b.disabled=following})}
function chooseAllSubscriptions(){
 el('read-source').value='';selectedSource='';renderReading();showView('news');
 el('news-content').scrollIntoView({block:'start',behavior:'smooth'});
}
function chooseSource(url){
 const source=feeds.find(f=>f.url===url);
 if(source&&excludedURL(publisherForFeed(source))){chooseAllSubscriptions();return}
 selectedSource=url;el('read-source').value='personal:'+url;renderReading();showView('news');
 el('reading-title').scrollIntoView({block:'start',behavior:'smooth'});
}
function renderSourceSelector(){
 const group=el('followed-options');group.replaceChildren();
 feeds.forEach(f=>{const option=document.createElement('option');option.value='personal:'+f.url;option.textContent=f.name;group.append(option)});
 if(selectedSource&&!feeds.some(f=>f.url===selectedSource))selectedSource='';
 const select=el('read-source');if(select.value.startsWith('personal:')&&!feeds.some(f=>'personal:'+f.url===select.value))select.value='';
}
// Every visible article is fetched from a subscribed feed. The server catalog
// is discovery metadata only; it never supplies unsolicited headlines.
function readingSelection(source){
 const all=source==='';
 const feedURL=source.startsWith('personal:')?source.slice(9):'';
 return {source,all,feedURL,showPersonal:true};
}
function renderReading(){
 // Also guard saved/selector-driven navigation to excluded publishers.
 if(selectedPublisherExcluded()){el('read-source').value='';selectedSource=''}
 const selection=readingSelection(el('read-source').value);
 selectedSource=selection.feedURL;
 el('personal-feed').classList.toggle('personal-hidden',!selection.showPersonal);
 el('subscriptions-empty-notice').hidden=!selection.all||registry.length>0;
 el('all-subscriptions').hidden=selection.all;
 el('reading-title').hidden=selection.all;
 el('reading-title').textContent=selection.feedURL?(feeds.find(f=>f.url===selection.feedURL)?.name||'Selected subscription'):'';
 el('personal-section-label').hidden=true;
 renderPersonalArticles();
}
let allPersonalArticles=[], personalErrors=[];
function renderPersonalArticles(){
 const lead=el('personal-lead-stories'),list=el('personal-articles');
 lead.replaceChildren();list.replaceChildren();
 const sources=new Set(feeds.filter(f=>!selectedSource||f.url===selectedSource).map(f=>f.name));
 const articles=allPersonalArticles.filter(a=>sources.has(a.source)&&!excludedURL(a.url));
 const seen=new Set();articles.sort((a,b)=>new Date(b.published)-new Date(a.published));
 let count=0;
 articles.forEach(a=>{
  if(seen.has(a.url))return;seen.add(a.url);
  const featured=count++<3;
  const row=document.createElement('article');row.className=featured?'feature':'item';
  const body=document.createElement('div');body.className='item-body';
  const meta=document.createElement('div');meta.className='source';meta.textContent=a.source+' · '+a.category;
  const title=document.createElement('a');title.className='story-title';title.href=a.open;title.target='newsgo-scrapergo';title.textContent=a.title;
  const sub=document.createElement('div');sub.className='story-meta';sub.textContent=a.date+'  ';
  const original=document.createElement('a');original.href=a.url;original.target='_blank';original.rel='noopener noreferrer';original.textContent=featured?'Original publisher ↗':'Original ↗';original.className='original';
  sub.append(original);body.append(meta,title,sub);const exclude=excludeButton(a.url);if(exclude)body.append(exclude);
  if(a.image&&/^https:\/\/[^\s/]+/i.test(a.image)){
   const img=document.createElement('img');img.className='article-thumb';img.src=a.image;img.alt='';img.loading='eager';img.referrerPolicy='no-referrer';img.addEventListener('error',()=>img.remove());
   if(featured)body.prepend(img);else row.append(img);
  }
  if(featured){row.append(body);lead.append(row)}else{row.prepend(body);list.append(row)}
 });
 const hasArticles=count>0;
 el('personal-lead-section').hidden=!hasArticles;
 el('personal-latest-section').hidden=!hasArticles;
 const selected=feeds.find(f=>f.url===selectedSource);
 if(!el('personal-loading').classList.contains('is-refreshing'))
  el('personal-loading').textContent=seen.size+' headlines'+(selected?' from '+selected.name:' from your subscriptions')+(personalErrors.length?' · '+personalErrors.join('; '):'');
}

function renderManager(){
 renderSourceSelector();
 const total=registry.length;el('personal-count').textContent='('+total+')';el('subscription-count').textContent='('+total+')';
 el('subscriptions-empty').textContent=registry.length?'':'No subscriptions yet. Choose Find Subscriptions to add one.';
 const root=el('personal-sources');root.replaceChildren();
 const allRow=document.createElement('div');allRow.className='personal-row';
 const allInfo=document.createElement('strong');allInfo.textContent='All Subscriptions';
 allRow.append(allInfo,button('Read All Articles',chooseAllSubscriptions));root.append(allRow);
 // Legacy subscription records that could not be mapped remain visible
 // with a clear explanation and an unsubscribe option; do not lose them.
 included.forEach(name=>{const row=document.createElement('div');row.className='personal-row';const info=document.createElement('strong');info.textContent=name+' · Feed unavailable — select a feed in Find Subscriptions';row.append(info,button('Unsubscribe',()=>{if(unsubscribeRecord('shared:'+name)){renderManager();refresh()}}));root.append(row)});
 document.querySelectorAll('.read-included').forEach(b=>{
  const followed=feeds.some(f=>f.url===b.dataset.feedUrl);
  b.textContent=followed?'Subscribed':'Subscribe';b.disabled=followed;
  const box=b.closest('.included-card');if(!box)return;
  const domain=exclusionDomain(b.dataset.feedUrl);
  box.dataset.excluded=excludedURL(domain)?'true':'false';
  let badge=box.querySelector('.exclusion-label');if(!badge){badge=document.createElement('small');badge.className='exclusion-label';box.insertBefore(badge,b)}
  badge.textContent=statusLabel(domain);
  let action=box.querySelector('.exclude-site');if(!action){action=excludeButton(domain);if(action)box.append(action)}
  if(action)action.textContent=excludedURL(domain)?'Allow this site':'Exclude this site';
 });
 feeds.forEach((f,i)=>{const row=document.createElement('div');row.className='personal-row';const info=document.createElement('div');const name=document.createElement('strong');name.textContent=f.name;const detail=document.createElement('small');try{detail.textContent=new URL(f.url).hostname+' · '+f.category}catch(e){detail.textContent=f.category}info.append(name,detail);const domain=publisherForFeed(f);markSite(info,domain);const actions=document.createElement('div');actions.className='subscription-actions';if(!domain||!excludedURL(domain))actions.append(button('Read Articles',()=>chooseSource(f.url)));actions.append(button('Unsubscribe',()=>{if(!unsubscribeRecord(feedIdentity(f.url)))return;
 if(selectedSource===f.url)selectedSource='';renderManager();refresh();message('Unsubscribed from '+f.name+'.')}));const exclude=excludeButton(domain);if(exclude)actions.append(exclude);row.append(info,actions);root.append(row)});
 syncFollowButtons();renderReading();
}
function showView(view){
 el('personal-panel').classList.toggle('personal-hidden',view!=='discover');
 el('subscriptions-panel').classList.toggle('personal-hidden',view!=='subscriptions');
 el('exclusions-panel').classList.toggle('personal-hidden',view!=='exclusions');
 el('news-content').classList.toggle('personal-hidden',view!=='news');
 [['manage-personal','discover'],['view-subscriptions','subscriptions'],['manage-exclusions','exclusions']].forEach(([id,value])=>el(id).setAttribute('aria-pressed',String(view===value)));
}
let refreshSequence=0;
async function refresh(){
 const sequence=++refreshSequence;
 const loading=el('personal-loading');
 renderReading();
 if(!feeds.length){loading.classList.remove('is-refreshing');allPersonalArticles=[];personalErrors=[];renderPersonalArticles();return}
 const snapshot=feeds.slice();
 loading.classList.add('is-refreshing');
 const context='';
 loading.textContent='Refreshing '+snapshot.length+' subscribed web feed'+(snapshot.length===1?'':'s')+'…'+context;
 try{
  let data=[];
  let completed=0;
  for(let offset=0;offset<snapshot.length;offset+=8){
   const batch=snapshot.slice(offset,offset+8);
   const res=await fetch('/api/personal',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({feeds:batch})});
   if(sequence!==refreshSequence)return;
   if(!res.ok)throw new Error('HTTP '+res.status);
   data.push(...await res.json());
   if(sequence!==refreshSequence)return;
   completed+=batch.length;
   if(completed<snapshot.length)loading.textContent='Updated '+completed+' of '+snapshot.length+' web feeds…'+context;
  }
  if(sequence!==refreshSequence)return;
  allPersonalArticles=[];personalErrors=[];
  data.forEach(group=>{if(group.error)personalErrors.push(group.source.name+': '+group.error);else allPersonalArticles.push(...(group.articles||[]))});
  loading.classList.remove('is-refreshing');
  renderPersonalArticles();
  const succeeded=data.filter(group=>!group.error).length;
  el('subscription-updated').textContent='Updated subscribed feeds';
  loading.textContent='Updated '+succeeded+' of '+snapshot.length+' web feeds'+(personalErrors.length?' · '+personalErrors.length+' failed':'')+context+' · '+allPersonalArticles.filter(a=>!excludedURL(a.url)).length+' visible headlines';
 }catch(e){
  if(sequence!==refreshSequence)return;
  loading.classList.remove('is-refreshing');
  loading.textContent='Could not refresh feeds: '+e.message;
 }
}

function add(f,trigger){if(!valid(f)){message('Enter a source name and valid HTTP(S) feed URL.');return}if(feeds.some(x=>x.url===f.url)){message('That feed is already in your personal list.');return}
 if(trigger){trigger.disabled=true;trigger.textContent='Checking…'}
 message('Validating feed…');
 fetch('/api/personal',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({feeds:[f]})}).then(async res=>{if(!res.ok){const detail=await res.text();throw Error('NewsGo validation returned HTTP '+res.status+(detail?': '+detail.slice(0,160):''))}const data=await res.json();if(data[0].error)throw Error(/HTTP 5\d\d/.test(data[0].error)?'The publisher feed server failed ('+data[0].error+'). Try another feed result or use Add a website.':data[0].error);if(!subscribeRecord(feedRecord(f)))return;renderManager();message('Added '+f.name+' — '+(data[0].articles||[]).length+' articles found.');refresh()}).catch(err=>{message('Feed not added: '+err.message);syncFollowButtons();if(trigger){trigger.disabled=false;trigger.textContent='Subscribe'}})}
async function discoverFeeds(){const u=el('discover-url').value.trim(),out=el('discover-results');out.replaceChildren();message('Searching for RSS or Atom feeds…');
 try{const r=await fetch('/api/discover',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({url:u})});const data=await r.json();if(!r.ok)throw Error(data.error||'HTTP '+r.status);if(!data.feeds||!data.feeds.length)throw Error('No RSS or Atom feeds found. Try Advanced with a direct feed URL.');
 data.feeds.forEach(f=>{const row=document.createElement('div');row.className='personal-row';const info=document.createElement('div');const title=document.createElement('strong');title.textContent=f.name;const small=document.createElement('small');small.textContent=(data.counts&&data.counts[f.url]||0)+' articles · '+(f.category||'News')+' · '+(new URL(f.url)).hostname;const details=document.createElement('details');const summary=document.createElement('summary');summary.textContent='Technical feed address';const address=document.createElement('small');address.textContent=f.url;details.append(summary,address);info.append(title,document.createElement('br'),small,details);const domain=publisherForFeed(f);markSite(info,domain);row.append(info,button('Subscribe',()=>add(f)));const exclusion=excludeButton(domain);if(exclusion)row.append(exclusion);out.append(row)});message(data.feeds.length+' feed(s) found. Select Subscribe to add one.');
 }catch(err){message('Discovery: '+err.message)}}
el('discover-button').addEventListener('click',discoverFeeds);
el('discover-form').addEventListener('submit',e=>{e.preventDefault();discoverFeeds()});
el('catalog-search').addEventListener('input',queueSearch);
el('manage-personal').addEventListener('click',()=>showView('discover'));

function chooseIncluded(name){const f=catalogByName.get(name);if(f)chooseSource(f.url)}
document.querySelectorAll('.read-included').forEach(b=>b.addEventListener('click',()=>{
 const f=catalogByName.get(b.dataset.includedSource);
 if(f)add(f,b);
}));
el('view-subscriptions').addEventListener('click',()=>showView('subscriptions'));
el('all-subscriptions').addEventListener('click',chooseAllSubscriptions);
el('close-discovery').addEventListener('click',()=>showView('news'));
el('close-subscriptions').addEventListener('click',()=>showView('news'));
el('personal-form').addEventListener('submit',e=>{e.preventDefault();add({name:el('personal-name').value.trim(),url:el('personal-url').value.trim(),category:el('personal-category').value});});
const sug=el('personal-suggestions');suggestions.forEach(s=>sug.append(button('+ '+s.name,()=>add(s),'personal-suggestion')));
el('personal-export').addEventListener('click',()=>{const blob=new Blob([JSON.stringify({format:'newsgo-preferences-v3',subscriptions:registry,excludedDomains},null,2)],{type:'application/json'}),u=URL.createObjectURL(blob),a=document.createElement('a');a.href=u;a.download='newsgo-preferences-v3.json';a.click();setTimeout(()=>URL.revokeObjectURL(u),1000)});
el('personal-import').addEventListener('click',()=>el('personal-file').click());
el('personal-file').addEventListener('change',async e=>{const f=e.target.files[0];if(!f)return;try{if(f.size>16384)throw Error('File too large');const data=JSON.parse(await f.text());const records=Array.isArray(data)?data:data&&data.format==='newsgo-preferences-v3'?data.subscriptions:null;if(!Array.isArray(records))throw Error('Invalid subscription list');
   const importedExclusions=Array.isArray(data)?null:data.excludedDomains;if(importedExclusions!==null&&(!Array.isArray(importedExclusions)||importedExclusions.some(d=>typeof d!=='string'||exclusionDomain(d)!==d)))throw Error('Invalid exclusion list');
   const migrated=convertLegacy(records.map(r=>r&&r.kind?r:feedRecord(r)));
 if(migrated.some(r=>!recordValid(r)))throw Error('Invalid subscription entry');
 const unique=[...new Map(migrated.map(r=>[r.id,r])).values()];
 if(!saveSubscriptions(unique))return;
   if(importedExclusions!==null&&!saveExclusions(importedExclusions))return;
   renderManager();refresh();message('Subscriptions imported'+(importedExclusions!==null?' with exclusions':' (existing exclusions preserved)')+'.')}catch(err){message('Import failed: '+err.message)}e.target.value=''});
el('exclude-add').addEventListener('click',()=>{excludeDomain(el('exclude-domain').value);el('exclude-domain').value=''});
el('exclude-domain').addEventListener('keydown',e=>{if(e.key==='Enter'){e.preventDefault();el('exclude-add').click()}});
el('manage-exclusions').addEventListener('click',()=>showView('exclusions'));
el('close-exclusions').addEventListener('click',()=>showView('news'));
message(storageError||exclusionError||'Search by topic or publisher to discover subscriptions.');renderExclusionManager();renderManager();showView('news');searchCatalog();refresh();setInterval(refresh,15*60*1000);

})();`
