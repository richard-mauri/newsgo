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
let excludedDomains=[],articleDomains=[],exclusionError='';
const articleExclusionKey='newsgo.excludedArticleDomains.v1';
const agePreferenceKey='newsgo.maxArticleAgeDays.v1';
const allowedAges=[0,1,3,7,14,30,60];
function normalizedAge(v){const n=Number(v);return allowedAges.includes(n)?n:0}
let maxArticleAgeDays=normalizedAge(localStorage.getItem(agePreferenceKey)||0);
function setArticleAge(days){maxArticleAgeDays=normalizedAge(days);localStorage.setItem(agePreferenceKey,String(maxArticleAgeDays));el('maximum-article-age').value=String(maxArticleAgeDays);scheduleServerSave();refresh();}

function hostOf(value){
 try{const u=new URL(value.includes('://')?value:'https://'+value);if(!/^https?:$/.test(u.protocol)||u.username||u.password)return '';return u.hostname.toLowerCase().replace(/\.$/,'').replace(/^www\./,'')}
 catch(_){return ''}
}
function exclusionDomain(value){
 const host=hostOf(value);
 return host&&/^[a-z0-9.-]+$/.test(host)&&host.includes('.')&&!host.startsWith('.')&&!host.endsWith('.')&&!host.includes('..')?host:'';
}
// NewsGo's publisher identity groups agency subdomains under a gov.au agency.
 // Other domains retain their existing exact-host behavior; don't guess public suffixes.
function publisherScope(value){
 const d=exclusionDomain(value);
 if(!d)return '';
 const parts=d.split('.');
 return d.endsWith('.gov.au')&&parts.length>3?parts.slice(-3).join('.'):d;
}
function publisherExcluded(value){
 const scope=publisherScope(value);
 return !!scope&&excludedDomains.some(d=>{
  const saved=publisherScope(d);
  return scope===saved||scope.endsWith('.'+saved);
 });
}
function articleDomain(value){const d=exclusionDomain(value);return d==='youtu.be'?'youtube.com':d}
function effectiveExcludedDomains(){return [...new Set([...excludedDomains,...articleDomains].map(publisherScope).filter(Boolean))]}
function domainIsExcluded(value){const scope=publisherScope(value);return !!scope&&effectiveExcludedDomains().some(d=>scope===d||scope.endsWith('.'+d))}
function articleExcluded(value){return domainIsExcluded(value)||(articleDomains.includes('youtube.com')&&hostOf(value)==='youtu.be')}
function confirmDomainRemoval(domain,affected){
 if(!affected.length)return Promise.resolve(true);
 return new Promise(resolve=>{
  const previous=document.activeElement;
  const shade=document.createElement('div');shade.className='sync-choice-backdrop';
  const dialog=document.createElement('section');dialog.className='sync-choice-dialog';dialog.setAttribute('role','dialog');dialog.setAttribute('aria-modal','true');dialog.setAttribute('aria-labelledby','exclude-confirm-title');
  const title=document.createElement('h3');title.id='exclude-confirm-title';title.textContent='Exclude '+domain+'?';
  const explanation=document.createElement('p');explanation.textContent='This will remove '+affected.length+' subscription(s) from My Subscriptions and block article links from this domain and its subdomains:';
  const list=document.createElement('ul');list.className='exclusion-confirm-list';
  affected.forEach(r=>{const item=document.createElement('li');item.textContent=r.name;list.append(item)});
  const note=document.createElement('p');note.textContent='Allowing the domain later will not restore removed subscriptions.';
  const actions=document.createElement('div');actions.className='sync-choice-actions';
  function finish(ok){document.removeEventListener('keydown',onKey);shade.remove();if(previous&&previous.isConnected)previous.focus();resolve(ok)}
  function onKey(e){if(e.key==='Escape'){e.preventDefault();finish(false)}}
  const cancel=button('Cancel',()=>finish(false));const confirm=button('Exclude domain and remove subscriptions',()=>finish(true));
  confirm.className='exclusion-confirm-action';
  actions.append(cancel,confirm);dialog.append(title,explanation,list,note,actions);shade.append(dialog);document.body.append(shade);
  document.addEventListener('keydown',onKey);cancel.focus();
 });
}
async function setUnifiedDomain(domain,allow){
 const nextPublishers=excludedDomains.filter(d=>publisherScope(d)!==domain);
 const nextArticles=articleDomains.filter(d=>publisherScope(d)!==domain);
 const affected=allow?[]:registry.filter(r=>r.kind==='feed'&&(publisherScope(r.url)===domain||publisherScope(r.url).endsWith('.'+domain)));
 if(!allow&&!(await confirmDomainRemoval(domain,affected)))return false;
 const nextRegistry=allow?registry:registry.filter(r=>!affected.includes(r));
 const prev=[localStorage.getItem(exclusionKey),localStorage.getItem(articleExclusionKey),localStorage.getItem(unifiedKey)];
 try{
  localStorage.setItem(exclusionKey,JSON.stringify(allow?nextPublishers:[...nextPublishers,domain]));
  localStorage.setItem(articleExclusionKey,JSON.stringify(allow?nextArticles:[...nextArticles,domain]));
  localStorage.setItem(unifiedKey,JSON.stringify(nextRegistry));
 }catch(e){
  [exclusionKey,articleExclusionKey,unifiedKey].forEach((k,i)=>{try{if(prev[i]===null)localStorage.removeItem(k);else localStorage.setItem(k,prev[i])}catch(_){}});
  message('Exclusion not saved; preferences left unchanged: '+e.message);return false;
 }
 excludedDomains=allow?nextPublishers:[...nextPublishers,domain];articleDomains=allow?nextArticles:[...nextArticles,domain];
 registry=nextRegistry;syncDerived();
 refreshExclusionViews();renderPersonalArticles();scheduleServerSave();refresh();
 message(allow?'Allowed '+domain:'Excluded '+domain+(affected.length?' and removed '+affected.length+' subscription(s)':''));
 return true;
}
function readArticleExclusions(){try{const v=JSON.parse(localStorage.getItem(articleExclusionKey)||'[]');if(!Array.isArray(v)||v.some(x=>typeof x!=='string'||articleDomain(x)!==x))throw Error('Invalid article-domain list');articleDomains=[...new Set(v)]}catch(e){exclusionError='Article exclusions were not changed: '+e.message}}
function saveArticleExclusions(next){try{const domains=[...new Set(next)];localStorage.setItem(articleExclusionKey,JSON.stringify(domains));articleDomains=domains;renderExclusionManager();renderPersonalArticles();scheduleServerSave();return true}catch(e){message('Could not save article exclusions: '+e.message);return false}}
function toggleArticleDomain(value){toggleDomain(value)}
function excludedURL(value){return domainIsExcluded(value)}
function readExclusions(){try{const data=JSON.parse(localStorage.getItem(exclusionKey)||'[]');if(!Array.isArray(data)||data.some(x=>typeof x!=='string'||exclusionDomain(x)!==x))throw Error('Invalid excluded-domain list');excludedDomains=[...new Set(data)]}catch(e){exclusionError='Exclusions were not changed: '+e.message}}
function saveExclusions(list){
 try{const unique=[...new Set(list)];localStorage.setItem(exclusionKey,JSON.stringify(unique));excludedDomains=unique;refreshExclusionViews();scheduleServerSave();return true}
 catch(e){message('Could not save exclusions: '+e.message);return false}
}
function excludeDomain(value){const domain=publisherScope(value);if(!domain){message('Enter a valid website domain.');return}if(!excludedURL(domain))setUnifiedDomain(domain,false);else message(domain+' is already excluded.')}
function toggleDomain(value){const domain=publisherScope(value);if(!domain){message('Could not identify a website domain.');return}setUnifiedDomain(domain,excludedURL(domain))}
function publisherForFeed(feed){return publisherScope(feed.website||feed.url||'')}
function statusLabel(domain){return domain&&excludedURL(domain)?' · Excluded':''}
function excludeButton(domain){const d=publisherScope(domain);if(!d)return null;const b=button(excludedURL(d)?'Allow domain':'Exclude This Domain',()=>toggleDomain(d),'exclude-site');b.dataset.publisherDomain=d;return b}
function markSite(info,domain){if(!domain)return;const badge=document.createElement('small');badge.className='exclusion-label';badge.textContent=statusLabel(domain);info.append(badge)}
// Compare against the specific excluded entry, not against every exclusion.
function excludedURLForDomain(url,domain){const host=hostOf(url);return !!host&&(host===domain||host.endsWith('.'+domain))}
function renderExclusionManager(){
 const root=el('excluded-list');root.replaceChildren();
 const domains=effectiveExcludedDomains().sort();
 el('excluded-count').textContent='('+domains.length+')';
 if(!domains.length){const p=document.createElement('p');p.textContent='No domains excluded.';root.append(p)}
 domains.forEach(d=>{
  const row=document.createElement('div');row.className='personal-row';
  const info=document.createElement('strong');info.textContent=d;
  row.append(info,button('Allow domain',()=>setUnifiedDomain(d,true)));root.append(row);
 });
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
 renderExclusionManager();renderManager();refreshSearchExclusionCards();
}
function refreshSearchExclusionCards(){
 document.querySelectorAll('#catalog-list .catalog-card').forEach(card=>{
  const action=card.querySelector('button.exclude-site');
  if(!action)return;
  const domain=action.dataset.publisherDomain;
  const excluded=excludedURL(domain);
  action.textContent=excluded?'Allow domain':'Exclude This Domain';
  card.dataset.excluded=excluded?'true':'false';
  let badge=card.querySelector('.exclusion-label');
  if(!badge){badge=document.createElement('small');badge.className='exclusion-label';card.querySelector('div')?.append(badge)}
  if(badge)badge.textContent=statusLabel(domain);
 });
 syncFollowButtons();
}
readExclusions();readArticleExclusions();
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
  registry=next;syncDerived();scheduleServerSave();return true;
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
const readingSelectionKey='newsgo.lastReadingSelection.v1';
let activeReadingSelection=null; // null until user chooses a feed, All, or restores a saved choice
function validReadingSelection(value){
 if(value==='all')return 'all';
 return typeof value==='string'&&feeds.some(f=>f.url===value&&!excludedURL(publisherForFeed(f)))?value:null;
}
function restoreReadingSelection(value){
 const selection=validReadingSelection(value);
 activeReadingSelection=selection;
 selectedSource=selection&&selection!=='all'?selection:'';
 el('read-source').value=selectedSource?'personal:'+selectedSource:'';
 renderReading();
 if(selection!==null){showView('news');refresh()}
 else{allPersonalArticles=[];personalErrors=[];filteredEntries=[];notProcessedEntries=[];renderPersonalArticles()}
}
function rememberReadingSelection(value){
 localStorage.setItem(readingSelectionKey,value);
 scheduleServerSave();
}

function message(t){el('personal-message').textContent=t}
function button(label,handler,cls){const b=document.createElement('button');b.type='button';b.textContent=label;if(cls)b.className=cls;b.addEventListener('click',handler);return b}
let searchSeq=0, lastSearchSentAt=0, searchInFlight=false;
const searchIntervalMs=3200; // Server allows one uncached search per IP every 3 seconds.
const searchDelay=ms=>new Promise(resolve=>setTimeout(resolve,ms));
async function searchCatalog(){
 const q=el('catalog-search').value.trim(),root=el('catalog-list'),status=el('catalog-message');
 if(searchInFlight)return;
 const seq=++searchSeq;
 if(q.length<2){root.replaceChildren();status.textContent='Enter at least two characters to search for news sources.';return}
 searchInFlight=true;el('catalog-submit').disabled=true;
 // Keep the previous results visible while waiting, but label them as previous results.
 status.setAttribute('role','status');status.setAttribute('aria-live','polite');
 status.textContent='⏳ Searching subscriptions for “'+q+'”… Previous results remain visible until updated.';
 try{
  let res,data;
  for(let attempt=0;attempt<3;attempt++){
   const wait=Math.max(0,searchIntervalMs-(Date.now()-lastSearchSentAt));
   if(wait>0){status.textContent='⏳ Waiting to continue search for “'+q+'”…';await searchDelay(wait)}
   if(seq!==searchSeq)return;
   lastSearchSentAt=Date.now();
   res=await fetch('/api/search-feeds?q='+encodeURIComponent(q),{cache:'no-store'});
   if(seq!==searchSeq)return;
   data=await res.json();
   if(res.status!==429)break;
   if(attempt===2)throw Error('Search is temporarily busy. Please try again.');
   status.textContent='⏳ Search rate limited; retrying automatically…';
  }
  if(seq!==searchSeq)return;
  if(!res.ok)throw Error(data.error||'Search unavailable (HTTP '+res.status+').');
  // Only replace the previous results after the new response was successful.
  root.replaceChildren();
  if(!data.feeds||!data.feeds.length){status.textContent='No matching sources found. Try another search or Add a website below.';return}
  data.feeds.forEach(f=>{
   const row=document.createElement('div');row.className='catalog-card';
   const info=document.createElement('div'),name=document.createElement('strong'),detail=document.createElement('small');
   name.textContent=f.name;detail.textContent=f.website||new URL(f.url).hostname;
   info.append(name,detail);const domain=publisherForFeed(f);markSite(info,domain);
   const followed=feeds.some(x=>x.url===f.url);
   const b=button(followed?'Subscribed':'Subscribe',()=>add({name:f.name,category:'Other',url:f.url},b));b.dataset.feedUrl=f.url;b.dataset.publisherDomain=domain;b.disabled=followed||excludedURL(domain);
   const exclude=excludeButton(domain);row.append(info,b);if(exclude)row.append(exclude);root.append(row);
  });
  syncFollowButtons();
  status.textContent=data.feeds.length+' candidate sources. Each feed is validated when you select Subscribe.';
 }catch(err){if(seq===searchSeq)status.textContent='Search failed: '+err.message+' Previous results remain visible; try again shortly.'}
 finally{searchInFlight=false;el('catalog-submit').disabled=false}
}
function queryEdited(){
 ++searchSeq; // Discard a response if its submitted query is no longer current.
 const q=el('catalog-search').value.trim();
 el('catalog-message').textContent=q.length>=2?'Press Search or Enter to search for “'+q+'”. Previous results stay visible.':'Enter at least two characters to search for news sources.';
}
function submitCatalogSearch(e){if(e)e.preventDefault();searchCatalog()}

function syncFollowButtons(){document.querySelectorAll('[data-feed-url]').forEach(b=>{const following=feeds.some(f=>f.url===b.dataset.feedUrl);const domain=b.dataset.publisherDomain||exclusionDomain(b.dataset.feedUrl);b.textContent=following?'Subscribed':'Subscribe';b.disabled=following||excludedURL(domain);if(!following&&excludedURL(domain))b.title='Allow this site before subscribing';else b.removeAttribute('title')})}
function chooseAllSubscriptions(){
 activeReadingSelection='all';rememberReadingSelection('all');el('read-source').value='';selectedSource='';renderReading();showView('news');refresh();
 el('news-content').scrollIntoView({block:'start',behavior:'smooth'});
}
function chooseSource(url){
 const source=feeds.find(f=>f.url===url);
 if(source&&excludedURL(publisherForFeed(source))){chooseAllSubscriptions();return}
 activeReadingSelection=url;rememberReadingSelection(url);selectedSource=url;el('read-source').value='personal:'+url;renderReading();showView('news');refresh();
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
 el('reading-title').hidden=true; // Single selected-feed heading is rendered above filtered results.
 el('reading-title').textContent=selection.feedURL?(feeds.find(f=>f.url===selection.feedURL)?.name||'Selected subscription'):'';
 el('personal-section-label').hidden=true;
 renderFilteredSubscriptionOffer();
 renderPersonalArticles();
}
let allPersonalArticles=[], personalErrors=[], filteredEntries=[], notProcessedEntries=[];

// This recommendation is based only on classification decisions actually made
// for the selected feed, never the All subscriptions view.
function filteredSubscriptionStats(sourceName) {
 const rejected=filteredEntries.filter(a=>a.source===sourceName&&
  !['too-old','publisher exclusion','article-domain exclusion'].includes(a.reason)).length;
 const accepted=allPersonalArticles.filter(a=>a.source===sourceName).length;
 const evaluated=rejected+accepted;
 return {rejected,accepted,evaluated,offer:evaluated>=10&&rejected/evaluated>0.9};
}
function confirmSubscriptionRemoval(name){
 return new Promise(resolve=>{
  const previous=document.activeElement;
  const shade=document.createElement('div');shade.className='sync-choice-backdrop';
  const dialog=document.createElement('section');dialog.className='sync-choice-dialog';
  dialog.setAttribute('role','dialog');dialog.setAttribute('aria-modal','true');dialog.setAttribute('aria-labelledby','unsubscribe-confirm-title');
  const title=document.createElement('h3');title.id='unsubscribe-confirm-title';title.textContent='Unsubscribe from '+name+'?';
  const description=document.createElement('p');description.textContent='This removes the subscription from My Subscriptions. It does not exclude the publisher or any article domain.';
  const actions=document.createElement('div');actions.className='sync-choice-actions';
  function finish(ok){document.removeEventListener('keydown',onKey);shade.remove();if(previous&&previous.isConnected)previous.focus();resolve(ok)}
  function onKey(event){if(event.key==='Escape'){event.preventDefault();finish(false)}}
  const cancel=button('Cancel',()=>finish(false));
  const confirm=button('Unsubscribe',()=>finish(true),'exclusion-confirm-action');
  actions.append(cancel,confirm);dialog.append(title,description,actions);shade.append(dialog);document.body.append(shade);
  document.addEventListener('keydown',onKey);cancel.focus();
 });
}
function renderFilteredSubscriptionOffer(){
 const anchor=el('filtered-articles-card');
 let heading=el('selected-subscription-heading');
 if(!heading){heading=document.createElement('h2');heading.id='selected-subscription-heading';anchor.insertAdjacentElement('beforebegin',heading);}
 heading.hidden=true;
 let panel=el('filtered-subscription-offer');
 if(!panel){
  panel=document.createElement('div');panel.id='filtered-subscription-offer';
  panel.className='personal-status';heading.insertAdjacentElement('afterend',panel);
 }
 panel.replaceChildren();panel.hidden=true;
 if(activeReadingSelection===null||activeReadingSelection==='all'||!selectedSource)return;
 const feed=feeds.find(f=>f.url===selectedSource);
 if(feed){heading.textContent=feed.name;heading.hidden=false;}
 if(!feed||personalErrors.some(e=>e.startsWith(feed.name+': '))||
  notProcessedEntries.some(e=>e.source===feed.name&&e.count>0))return;
 const stats=filteredSubscriptionStats(feed.name);
 if(!stats.offer)return;
 const pct=Math.round(100*stats.rejected/stats.evaluated);
 const message=document.createElement('span');
 message.textContent=pct+'% of '+stats.evaluated+' classified articles from '+feed.name+
  ' were filtered ('+stats.rejected+' filtered). ';
 const action=button('Unsubscribe from '+feed.name,async()=>{
  // Confirmation is required; displaying a recommendation never alters a subscription.
  if(!await confirmSubscriptionRemoval(feed.name))return;
  const record=registry.find(r=>r.kind==='feed'&&r.url===feed.url);
  if(!record||!unsubscribeRecord(record.id))return;
  chooseAllSubscriptions();
 });
 action.className='exclusion-confirm-action';
 action.style.marginLeft='0.75rem';
 panel.append(message,action);panel.hidden=false;
}

function renderFilteredArticles(){
 const card=el('filtered-articles-card'),list=el('filtered-articles-list'),count=el('filtered-articles-count');
 if(!card||!list)return;
 list.replaceChildren();
 const sources=new Set(feeds.filter(f=>activeReadingSelection!==null&&(!selectedSource||f.url===selectedSource)).map(f=>f.name));
 let display=[];
 for(const f of filteredEntries)if(sources.has(f.source))display.push(f);
 for(const a of allPersonalArticles)if(sources.has(a.source)){
   if(excludedURL(a.url))display.push({title:a.title,url:a.url,source:a.source,reason:'publisher exclusion'});
   else if(articleExcluded(a.url))display.push({title:a.title,url:a.url,source:a.source,reason:'article-domain exclusion'});
 }
 const pending=notProcessedEntries.filter(x=>sources.has(x.source));
 count.textContent='('+display.length+' filtered, '+pending.reduce((n,x)=>n+x.count,0)+' unprocessed)';
 card.hidden=activeReadingSelection===null;
 if(display.length===0&&pending.length===0){const empty=document.createElement('p');empty.className='personal-status';empty.textContent='No filtered or unprocessed entries reported for this selection.';list.append(empty);}
 for(const a of display){const row=document.createElement('p');row.className='personal-row';const link=document.createElement('a');link.href=a.url;link.target='_blank';link.rel='noopener noreferrer';link.textContent=a.title||a.url;const reason=document.createElement('small');reason.textContent=' '+a.source+' — '+(a.reason==='too-old'?'Older than '+maxArticleAgeDays+' days':(a.reason||'unknown'))+(a.published&&a.reason==='too-old'?' · Published '+new Date(a.published).toLocaleDateString():'');row.append(link,reason);
 if(syncConnected&&isAdmin&&a.url&&a.reason!=='too-old'&&a.reason!=='publisher exclusion'&&a.reason!=='article-domain exclusion'){
  const details=document.createElement('details');
  const summary=document.createElement('summary');summary.textContent='Admin diagnostics (run fresh classification)';
  const output=document.createElement('pre');output.textContent='Open to inspect a fresh classification request.';
  details.append(summary,output);
  details.addEventListener('toggle',async()=>{
   if(!details.open)return;
   output.textContent='Checking…';
   try{
    const response=await fetch('/api/admin/classification-diagnostic',{method:'POST',headers:{Authorization:'Bearer '+syncToken,'Content-Type':'application/json'},body:JSON.stringify({url:a.url}),cache:'no-store'});
    if(!response.ok)throw Error('HTTP '+response.status);
    const result=await response.json();
    output.textContent='Classification: '+(result.classification||'unknown')+'\\nFailure: '+(result.failure||'none')+'\\nHTTP status: '+(result.httpStatus||'not available')+'\\nElapsed: '+result.elapsedMs+' ms\\nTimeout: '+result.timeoutSeconds+' seconds';
   }catch(error){output.textContent='Diagnostic unavailable: '+error.message;}
  });
  row.append(details);
 }
 list.append(row)}
 for(const x of pending){const p=document.createElement('p');p.textContent=x.source+': '+x.count+' additional feed entries not processed (per-feed cap)';list.append(p)}
}

function renderPersonalArticles(){
 renderFilteredSubscriptionOffer();
 const lead=el('personal-lead-stories'),list=el('personal-articles');
 lead.replaceChildren();list.replaceChildren();
 const sources=new Set(feeds.filter(f=>activeReadingSelection!==null&&(!selectedSource||f.url===selectedSource)).map(f=>f.name));
 const articles=allPersonalArticles.filter(a=>sources.has(a.source)&&!articleExcluded(a.url));
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
  if(a.dateWarning){const warning=document.createElement('span');warning.className='date-warning';warning.textContent=' ⚠ '+a.dateWarning;warning.title='Original feed date shown. Article age could not be determined.';sub.append(warning);}
  const original=document.createElement('a');original.href=a.url;original.target='_blank';original.rel='noopener noreferrer';original.textContent=featured?'Original publisher ↗':'Original ↗';original.className='original';
  sub.append(original);body.append(meta,title,sub);const exclude=excludeButton(a.url);if(exclude)body.append(exclude);const ad=articleDomain(a.url);if(ad){body.append(button('Exclude articles from '+ad,()=>toggleArticleDomain(ad),'exclude-article-domain'))}
  if(a.image&&/^https:\/\/[^\s/]+/i.test(a.image)){
   const img=document.createElement('img');img.className='article-thumb';img.src=a.image;img.alt='';img.loading='eager';img.referrerPolicy='no-referrer';img.addEventListener('error',()=>img.remove());
   if(featured)body.prepend(img);else row.append(img);
  }
  if(featured){row.append(body);lead.append(row)}else{row.prepend(body);list.append(row)}
 });
 renderFilteredArticles();
 const hasArticles=count>0;
 el('personal-lead-section').hidden=!hasArticles;
 el('personal-latest-section').hidden=!hasArticles;
 const selected=feeds.find(f=>f.url===selectedSource);
 if(activeReadingSelection===null){el('personal-loading').textContent='Welcome to NewsGo. Find Subscriptions, or select one subscription or All Subscriptions to read.';return}
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
 if(selectedSource===f.url)selectedSource='';renderManager();refresh();message('Unsubscribed from '+f.name+'.')}));if(!excludedURL(domain)){const exclude=excludeButton(domain);if(exclude)actions.append(exclude)}row.append(info,actions);root.append(row)});
 syncFollowButtons();renderReading();
}
function showView(view){
 el('personal-panel').classList.toggle('personal-hidden',view!=='discover');
 el('subscriptions-panel').classList.toggle('personal-hidden',view!=='subscriptions');
 el('exclusions-panel').classList.toggle('personal-hidden',view!=='exclusions');
 el('news-content').classList.toggle('personal-hidden',view!=='news');
 el('sync-panel').classList.toggle('personal-hidden',view!=='sync');
 el('admin-panel').classList.toggle('personal-hidden',view!=='admin');
 el('preferences-panel').classList.toggle('personal-hidden',view!=='preferences');
 [['manage-personal','discover'],['view-subscriptions','subscriptions'],['manage-preferences','preferences'],['manage-exclusions','exclusions'],['manage-sync','sync'],['manage-admin','admin']].forEach(([id,value])=>el(id).setAttribute('aria-pressed',String(view===value)));
}
let refreshSequence=0;
async function refresh(){
 const sequence=++refreshSequence;
 const loading=el('personal-loading');
 if(activeReadingSelection===null)return;
 renderReading();
 if(!feeds.length){loading.classList.remove('is-refreshing');allPersonalArticles=[];personalErrors=[];renderPersonalArticles();return}
 const snapshot=activeReadingSelection==='all'?feeds.slice():feeds.filter(f=>f.url===activeReadingSelection);
 loading.classList.add('is-refreshing');
 const context='';
 loading.textContent='Refreshing '+snapshot.length+' subscribed web feed'+(snapshot.length===1?'':'s')+'…'+context;
 try{
  let data=[];
  let completed=0;
  for(let offset=0;offset<snapshot.length;offset+=8){
   const batch=snapshot.slice(offset,offset+8);
   const res=await fetch('/api/personal',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({feeds:batch,maxArticleAgeDays})});
   if(sequence!==refreshSequence)return;
   if(!res.ok)throw new Error('HTTP '+res.status);
   data.push(...await res.json());
   if(sequence!==refreshSequence)return;
   completed+=batch.length;
   if(completed<snapshot.length)loading.textContent='Updated '+completed+' of '+snapshot.length+' web feeds…'+context;
  }
  if(sequence!==refreshSequence)return;
  allPersonalArticles=[];personalErrors=[];filteredEntries=[];notProcessedEntries=[];
  data.forEach(group=>{if(group.filtered)filteredEntries.push(...group.filtered);if(group.notProcessed)notProcessedEntries.push({source:group.source.name,count:group.notProcessed})});
  data.forEach(group=>{if(group.error)personalErrors.push(group.source.name+': '+group.error);else allPersonalArticles.push(...(group.articles||[]))});
  loading.classList.remove('is-refreshing');
  renderPersonalArticles();
  const succeeded=data.filter(group=>!group.error).length;
  loading.textContent='Updated '+succeeded+' of '+snapshot.length+' web feeds'+(personalErrors.length?' · '+personalErrors.length+' failed':'')+context+' · '+allPersonalArticles.filter(a=>!excludedURL(a.url)&&!articleExcluded(a.url)).length+' visible headlines';
 }catch(e){
  if(sequence!==refreshSequence)return;
  loading.classList.remove('is-refreshing');
  loading.textContent='Could not refresh feeds: '+e.message;
 }
}

function add(f,trigger){if(!valid(f)){message('Enter a source name and valid HTTP(S) feed URL.');return}if(excludedURL(trigger?.dataset.publisherDomain||publisherForFeed(f))){message('Allow this site before subscribing.');syncFollowButtons();return}if(feeds.some(x=>x.url===f.url)){message('That feed is already in your personal list.');return}
 if(trigger){trigger.disabled=true;trigger.textContent='Checking…'}
 message('Validating feed…');
 fetch('/api/personal',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({feeds:[f],maxArticleAgeDays})}).then(async res=>{if(!res.ok){const detail=await res.text();throw Error('NewsGo validation returned HTTP '+res.status+(detail?': '+detail.slice(0,160):''))}const data=await res.json();if(data[0].error)throw Error(/HTTP 5\d\d/.test(data[0].error)?'The publisher feed server failed ('+data[0].error+'). Try another feed result or use Add a website.':data[0].error);if(excludedURL(trigger?.dataset.publisherDomain||publisherForFeed(f)))throw Error('This site was excluded while the subscription was being checked. Allow the site before subscribing.');if(!subscribeRecord(feedRecord(f)))return;renderManager();message('Added '+f.name+' — '+(data[0].articles||[]).length+' articles found.');refresh()}).catch(err=>{message('Feed not added: '+err.message);syncFollowButtons();if(trigger){trigger.disabled=false;trigger.textContent='Subscribe'}})}
async function discoverFeeds(){const u=el('discover-url').value.trim(),out=el('discover-results');out.replaceChildren();message('Searching for RSS or Atom feeds…');
 try{const r=await fetch('/api/discover',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({url:u})});const data=await r.json();if(!r.ok)throw Error(data.error||'HTTP '+r.status);if(!data.feeds||!data.feeds.length)throw Error('No RSS or Atom feeds found. Try Advanced with a direct feed URL.');
 data.feeds.forEach(f=>{const row=document.createElement('div');row.className='personal-row';const info=document.createElement('div');const title=document.createElement('strong');title.textContent=f.name;const small=document.createElement('small');small.textContent=(data.counts&&data.counts[f.url]||0)+' articles · '+(f.category||'News')+' · '+(new URL(f.url)).hostname;const details=document.createElement('details');const summary=document.createElement('summary');summary.textContent='Technical feed address';const address=document.createElement('small');address.textContent=f.url;details.append(summary,address);info.append(title,document.createElement('br'),small,details);const domain=publisherForFeed(f);markSite(info,domain);row.append(info,button('Subscribe',()=>add(f)));const exclusion=excludeButton(domain);if(exclusion)row.append(exclusion);out.append(row)});message(data.feeds.length+' feed(s) found. Select Subscribe to add one.');
 }catch(err){message('Discovery: '+err.message)}}
el('discover-button').addEventListener('click',discoverFeeds);
el('discover-form').addEventListener('submit',e=>{e.preventDefault();discoverFeeds()});
el('catalog-search').addEventListener('input',queryEdited);
el('catalog-form').addEventListener('submit',submitCatalogSearch);
el('manage-personal').addEventListener('click',()=>showView('discover'));

function chooseIncluded(name){const f=catalogByName.get(name);if(f)chooseSource(f.url)}
document.querySelectorAll('.read-included').forEach(b=>b.addEventListener('click',()=>{
 const f=catalogByName.get(b.dataset.includedSource);
 if(f)add(f,b);
}));
el('view-subscriptions').addEventListener('click',()=>showView('subscriptions'));
el('read-source').addEventListener('change',()=>{const value=el('read-source').value;if(value.startsWith('personal:'))chooseSource(value.slice(9));else chooseAllSubscriptions()});
el('close-discovery').addEventListener('click',()=>showView('news'));
el('close-subscriptions').addEventListener('click',()=>showView('news'));
el('personal-form').addEventListener('submit',e=>{e.preventDefault();add({name:el('personal-name').value.trim(),url:el('personal-url').value.trim(),category:el('personal-category').value});});
const sug=el('personal-suggestions');suggestions.forEach(s=>sug.append(button('+ '+s.name,()=>add(s),'personal-suggestion')));
el('personal-export').addEventListener('click',()=>{const blob=new Blob([JSON.stringify({format:'newsgo-preferences-v3',subscriptions:registry,excludedDomains,excludedArticleDomains:articleDomains},null,2)],{type:'application/json'}),u=URL.createObjectURL(blob),a=document.createElement('a');a.href=u;a.download='newsgo-preferences-v3.json';a.click();setTimeout(()=>URL.revokeObjectURL(u),1000)});
el('personal-import').addEventListener('click',()=>el('personal-file').click());
el('personal-file').addEventListener('change',async e=>{const f=e.target.files[0];if(!f)return;try{if(f.size>16384)throw Error('File too large');const data=JSON.parse(await f.text());const records=Array.isArray(data)?data:data&&data.format==='newsgo-preferences-v3'?data.subscriptions:null;if(!Array.isArray(records))throw Error('Invalid subscription list');
   const importedExclusions=Array.isArray(data)?null:data.excludedDomains;if(importedExclusions!==null&&(!Array.isArray(importedExclusions)||importedExclusions.some(d=>typeof d!=='string'||exclusionDomain(d)!==d)))throw Error('Invalid exclusion list');
   const importedArticles=Array.isArray(data)?null:(data.excludedArticleDomains===undefined?null:data.excludedArticleDomains);if(importedArticles!==null&&(!Array.isArray(importedArticles)||importedArticles.some(d=>typeof d!=='string'||articleDomain(d)!==d)))throw Error('Invalid article exclusion list');
   const migrated=convertLegacy(records.map(r=>r&&r.kind?r:feedRecord(r)));
 if(migrated.some(r=>!recordValid(r)))throw Error('Invalid subscription entry');
 const unique=[...new Map(migrated.map(r=>[r.id,r])).values()];
 if(!saveSubscriptions(unique))return;
   if(importedExclusions!==null&&!saveExclusions(importedExclusions))return;
   if(importedArticles!==null&&!saveArticleExclusions(importedArticles))return;
   renderManager();refresh();message('Subscriptions imported'+(importedExclusions!==null?' with exclusions':' (existing exclusions preserved)')+'.')}catch(err){message('Import failed: '+err.message)}e.target.value=''});
el('exclude-add').addEventListener('click',()=>{excludeDomain(el('exclude-domain').value);el('exclude-domain').value=''});
el('exclude-domain').addEventListener('keydown',e=>{if(e.key==='Enter'){e.preventDefault();el('exclude-add').click()}});
el('article-exclude-add').addEventListener('click',()=>{toggleArticleDomain(el('article-exclude-domain').value);el('article-exclude-domain').value=''});
el('article-exclude-domain').addEventListener('keydown',e=>{if(e.key==='Enter'){e.preventDefault();el('article-exclude-add').click()}});
el('manage-exclusions').addEventListener('click',()=>{renderExclusionManager();showView('exclusions')});
el('close-exclusions').addEventListener('click',()=>showView('news'));

// Optional per-user server sync. The same token works on every device.
// Remembered credentials are local to this browser origin; only use on a trusted device.
const syncCredentialKey='newsgo.sync.token.v1';
let isAdmin=false,adminRevision=0;
let syncToken='',syncRevision=0,syncConnected=false,syncPending=false,syncBusy=false,syncTimer=null,syncDirty=false,syncConflict=false,syncReading=false;
const syncStatus=t=>{el('sync-status').textContent=t};
function syncAppearance(){
 el('sync-connection').textContent=syncConnected?'Connected':'Not connected';
 el('sync-connect').hidden=syncConnected;
 el('sync-token').hidden=syncConnected;
 el('manage-admin').hidden=!(syncConnected&&isAdmin);
 if(!(syncConnected&&isAdmin)&&!el('admin-panel').classList.contains('personal-hidden'))showView('news');
}
function syncSnapshot(){return {revision:syncRevision,subscriptions:registry,excludedDomains:excludedDomains,excludedArticleDomains:articleDomains,maxArticleAgeDays,lastReadingSelection:validReadingSelection(activeReadingSelection)||''}}
async function syncRequest(method,body){
 const response=await fetch('/api/sync',{method,headers:{Authorization:'Bearer '+syncToken,...(body?{'Content-Type':'application/json'}:{})},...(body?{body:JSON.stringify(body)}:{}) ,cache:'no-store'});
 if(response.status===404){const data=await response.json();return {...data,notFound:true}}
 if(!response.ok)throw Error('HTTP '+response.status+(response.status===409?' — newer server preferences are available. See Advanced to resolve the conflict.':' — '+(await response.text()).trim().slice(0,140)));
 return response.json();
}
function applyServerPreferences(data){
 if(!Array.isArray(data.subscriptions)||!Array.isArray(data.excludedDomains)||!Array.isArray(data.excludedArticleDomains||[]))throw Error('Invalid server profile');
 if(data.subscriptions.some(r=>!recordValid(r))||data.excludedDomains.some(d=>typeof d!=='string'||exclusionDomain(d)!==d)||!Array.isArray(data.excludedArticleDomains||[])||(data.excludedArticleDomains||[]).some(d=>typeof d!=='string'||articleDomain(d)!==d))throw Error('Server profile contains invalid entries');
 const migrated=convertLegacy(data.subscriptions);
 localStorage.setItem(unifiedKey,JSON.stringify(migrated));localStorage.setItem(exclusionKey,JSON.stringify(data.excludedDomains));localStorage.setItem(articleExclusionKey,JSON.stringify(data.excludedArticleDomains||[]));
 registry=migrated;syncDerived();excludedDomains=[...new Set(data.excludedDomains)];articleDomains=[...new Set(data.excludedArticleDomains||[])];maxArticleAgeDays=normalizedAge(data.maxArticleAgeDays||0);localStorage.setItem(agePreferenceKey,String(maxArticleAgeDays));el('maximum-article-age').value=String(maxArticleAgeDays);syncRevision=data.revision;
 const restored=validReadingSelection(data.lastReadingSelection);if(restored===null)localStorage.removeItem(readingSelectionKey);else localStorage.setItem(readingSelectionKey,restored);
 restoreReadingSelection(restored);
 syncDirty=false;syncPending=false;syncConflict=false;
 refreshExclusionViews();renderExclusionManager();renderManager();refresh();
}
// Require an explicit choice before replacing browser-only edits at connection.
function samePreferences(data){
 const a=convertLegacy(data.subscriptions),b=convertLegacy(registry);
 const ids=x=>x.map(r=>JSON.stringify(r)).sort().join('\n');
 const domains=x=>[...new Set(x)].sort().join('\n');
 return ids(a)===ids(b)&&domains(data.excludedDomains)===domains(excludedDomains)&&domains(data.excludedArticleDomains||[])===domains(articleDomains)&&normalizedAge(data.maxArticleAgeDays||0)===maxArticleAgeDays&&(activeReadingSelection===null||validReadingSelection(data.lastReadingSelection)===validReadingSelection(activeReadingSelection));
}
function resolveConnectDifference(){
 return new Promise(resolve=>{
  const shade=document.createElement('div');shade.className='sync-choice-backdrop';
  shade.setAttribute('role','presentation');
  const dialog=document.createElement('section');dialog.className='sync-choice-dialog';
  dialog.setAttribute('role','dialog');dialog.setAttribute('aria-modal','true');dialog.setAttribute('aria-labelledby','sync-choice-title');
  const h=document.createElement('h3');h.id='sync-choice-title';h.textContent='Choose your subscription preferences';
  const p=document.createElement('p');p.textContent='This browser has different reading preferences, subscriptions, or exclusions than your saved server profile. Choose whether to load the server settings or merge saved subscriptions and exclusions while retaining the local article-age preference.';
  const actions=document.createElement('div');actions.className='sync-choice-actions';
  function add(label,value){const b=document.createElement('button');b.type='button';b.textContent=label;b.addEventListener('click',()=>{shade.remove();resolve(value)});actions.append(b);return b}
  const useServer=add('Use Server Preferences','server');
  add('Merge Preferences','merge');add('Cancel','cancel');
  dialog.append(h,p,actions);shade.append(dialog);document.body.append(shade);useServer.focus();
 });
}
async function connectSync(token=el('sync-token').value.trim(),remember=true){
 if(!/^[a-f0-9]{64}$/i.test(token)){syncStatus('Enter the 64-character approved token distributed by Richard.');return}
 syncToken=token;syncConnected=false;isAdmin=false;syncReading=true;syncAppearance();syncStatus('Connecting…');
 let needsSave=false;
 try{
  const data=await syncRequest('GET');
  isAdmin=data.role==='admin';
  if(data.notFound){syncRevision=0;syncConnected=true;syncDirty=true;needsSave=true;syncStatus('Connected · creating your server profile…')}
  else {
   if(!Array.isArray(data.subscriptions)||!Array.isArray(data.excludedDomains)||data.subscriptions.some(r=>!recordValid(r))||data.excludedDomains.some(d=>typeof d!=='string'||exclusionDomain(d)!==d)||!Array.isArray(data.excludedArticleDomains||[])||(data.excludedArticleDomains||[]).some(d=>typeof d!=='string'||articleDomain(d)!==d))throw Error('Invalid server profile');
   if(!samePreferences(data)){
    const decision=await resolveConnectDifference();
    if(decision==='cancel'){
     syncToken='';syncConnected=false;localStorage.removeItem(syncCredentialKey);syncStatus('Connection cancelled · browser preferences preserved.');syncAppearance();return;
    }
    if(decision==='merge'){
     const combined=[...new Map([...convertLegacy(data.subscriptions),...registry].map(r=>[r.id,r])).values()];
     const domains=[...new Set([...data.excludedDomains,...excludedDomains])];
     const articles=[...new Set([...(data.excludedArticleDomains||[]),...articleDomains])];
     // The local reading-age preference wins when merging profiles.
     // Only replace local storage once both lists are validated.
     localStorage.setItem(unifiedKey,JSON.stringify(combined));localStorage.setItem(exclusionKey,JSON.stringify(domains));localStorage.setItem(articleExclusionKey,JSON.stringify(articles));
     registry=combined;syncDerived();excludedDomains=domains;articleDomains=articles;syncRevision=data.revision;
     // Merge keeps an explicit local selection; otherwise adopts the server choice.
     const mergedSelection=validReadingSelection(activeReadingSelection)||validReadingSelection(data.lastReadingSelection);
     if(mergedSelection!==null)localStorage.setItem(readingSelectionKey,mergedSelection);
     restoreReadingSelection(mergedSelection);
     refreshExclusionViews();renderManager();refresh();
     syncConnected=true;syncDirty=true;needsSave=true;
     syncStatus('Connected · saving merged preferences…');
    }else{applyServerPreferences(data);syncConnected=true;syncStatus('Connected · server preferences loaded (revision '+syncRevision+')')}
   }else{syncRevision=data.revision;syncConnected=true;syncDirty=false;restoreReadingSelection(validReadingSelection(data.lastReadingSelection));syncStatus('Connected · up to date (revision '+syncRevision+')')}
  }
  if(remember){try{localStorage.setItem(syncCredentialKey,token)}catch(e){syncStatus('Connected, but browser cannot remember your token: '+e.message)}}
  el('sync-token').value='';syncAppearance();
  // Update counters and management list even when the server profile is already current.
  if(syncConnected)renderManager();
 }catch(e){syncToken='';syncConnected=false;isAdmin=false;syncStatus('Could not connect: '+e.message);syncAppearance();if(remember)localStorage.removeItem(syncCredentialKey)}
 finally{syncReading=false;if(needsSave&&syncConnected)scheduleServerSave()}
}
async function loadSync(force=false){
 if(!syncConnected||syncBusy||syncReading)return;
 if((syncDirty||syncConflict)&&!force){syncStatus('Local changes need attention. Check Advanced before loading.');return}
 syncReading=true;
 try{
  const data=await syncRequest('GET');
  if(data.notFound){if(syncRevision>0)syncStatus('Server profile missing. Your local preferences were preserved.');return}
  if(data.revision===syncRevision&&!force)return;
  if(syncDirty&&!force)return;
  applyServerPreferences(data);syncStatus('Connected · synchronized (revision '+syncRevision+')');
 }catch(e){syncStatus('Sync unavailable: '+e.message+' · local preferences kept')}
 finally{syncReading=false}
}
async function saveSync(){
 if(!syncConnected||syncConflict||syncReading)return;
 if(syncBusy){syncPending=true;return}
 syncBusy=true;syncPending=false;
 const snap=syncSnapshot();
 try{
  const data=await syncRequest('PUT',snap);syncRevision=data.revision;
  const changed=JSON.stringify(snap.subscriptions)!==JSON.stringify(registry)||JSON.stringify(snap.excludedDomains)!==JSON.stringify(excludedDomains)||JSON.stringify(snap.excludedArticleDomains)!==JSON.stringify(articleDomains)||snap.maxArticleAgeDays!==maxArticleAgeDays||snap.lastReadingSelection!==(validReadingSelection(activeReadingSelection)||'');
  syncDirty=changed;if(changed)syncPending=true;
  syncStatus(changed?'Connected · saving newer changes…':'Connected · all changes saved');
 }catch(e){
  syncDirty=true;
  if(/HTTP 409/.test(e.message))syncConflict=true;
  syncStatus('Not synced: '+e.message+' · your local changes were preserved.');
  syncBusy=false;
  if(!syncConflict){clearTimeout(syncTimer);syncTimer=setTimeout(saveSync,30000)}
  return;
 }
 syncBusy=false;if(syncPending){syncPending=false;scheduleServerSave()}
}
function scheduleServerSave(){
 if(!syncConnected)return;
 syncDirty=true;
 syncStatus('Connected · saving changes…');
 clearTimeout(syncTimer);
 if(!syncConflict)syncTimer=setTimeout(saveSync,600);
}
async function checkSync(){
 if(!syncConnected||syncBusy||syncReading||syncDirty||syncConflict||document.hidden)return;
 await loadSync();
}
el('manage-sync').addEventListener('click',()=>showView('sync'));
el('close-sync').addEventListener('click',()=>showView('news'));
el('sync-connect').addEventListener('click',()=>connectSync());
el('sync-load').addEventListener('click',()=>{
 if(syncDirty||syncConflict){if(!confirm('Replace this browser’s subscriptions and excluded sites with the server copy? Unsaved local changes will be lost.'))return}
 clearTimeout(syncTimer);loadSync(true);
});
el('sync-save').addEventListener('click',()=>{if(!syncConnected){syncStatus('Connect first.');return}if(syncConflict){syncStatus('Server changed. Use Load from Server to resolve before saving.');return}clearTimeout(syncTimer);syncDirty=true;saveSync()});
el('sync-disconnect').addEventListener('click',()=>{
 clearTimeout(syncTimer);syncConnected=false;syncToken='';isAdmin=false;syncRevision=0;syncDirty=false;syncConflict=false;
 localStorage.removeItem(syncCredentialKey);el('sync-token').value='';syncAppearance();syncStatus('Disconnected · browser preferences remain available.');
});
const articleAdminFields={maxArticlesPerFeed:'admin-articles',classificationTimeoutSeconds:'admin-classification-timeout',classificationConcurrency:'admin-classification-parallel',feedCacheMinutes:'admin-feed-cache',feedSearchTimeoutSeconds:'admin-feed-search-timeout'};
let retainedLegacyImageSettings=null;
let adminSavedValues=null;
const adminStatus=message=>{el('admin-status').textContent=message};
const adminValues=()=>Object.fromEntries(Object.entries(articleAdminFields).map(([key,id])=>[key,el(id).value]));
const adminDirty=()=>adminSavedValues!==null&&Object.keys(articleAdminFields).some(key=>adminValues()[key]!==adminSavedValues[key]);
async function adminRequest(method,payload){
 const response=await fetch('/api/admin/settings',{method,headers:{Authorization:'Bearer '+syncToken,...(payload?{'Content-Type':'application/json'}:{})},...(payload?{body:JSON.stringify(payload)}:{}),cache:'no-store'});
 if(!response.ok)throw Error('HTTP '+response.status+' — '+(await response.text()).trim().slice(0,140));return response.json();
}
async function loadAdmin(){
 if(!syncConnected||!isAdmin)return false;
 try{const data=await adminRequest('GET');adminRevision=data.revision;
  retainedLegacyImageSettings=data.image;
  for(const [key,id] of Object.entries(articleAdminFields)){const field=el(id);if(field)field.value=(data.articles||{})[key]??'';}
  adminSavedValues=adminValues();
  adminStatus('Server settings loaded (revision '+adminRevision+').');
  return true;
 }catch(e){adminStatus('Could not load settings: '+e.message);return false}
}
async function saveAdmin(){
 if(!syncConnected||!isAdmin||!retainedLegacyImageSettings)return false;
 const form=el('admin-form');
 if(!form.reportValidity())return false;
 const articles=Object.fromEntries(Object.entries(adminValues()).map(([key,value])=>[key,Number(value)]));
 try{const result=await adminRequest('PUT',{revision:adminRevision,image:retainedLegacyImageSettings,articles});
  adminRevision=result.revision;retainedLegacyImageSettings=result.image;adminSavedValues=adminValues();
  adminStatus('Saved server settings (revision '+adminRevision+').');return true;
 }catch(e){adminStatus('Save failed: '+e.message+' Reload before retrying if another administrator made changes.');return false}
}
function askAdminUnsaved(){
 return new Promise(resolve=>{
  const dialog=el('admin-unsaved-dialog');
  if(typeof dialog.showModal!=='function'){
   resolve(window.confirm('Discard unsaved administrator settings?')?'discard':'cancel');return;
  }
  const done=()=>{dialog.removeEventListener('close',onClose);resolve(dialog.returnValue||'cancel')};
  const onClose=()=>done();
  dialog.addEventListener('close',onClose,{once:true});
  dialog.showModal();
 });
}
async function confirmAdminLeave(action){
 if(!adminDirty()){await action();return}
 const choice=await askAdminUnsaved();
 if(choice==='cancel')return;
 if(choice==='save' && !(await saveAdmin()))return;
 await action();
}
el('maximum-article-age').value=String(maxArticleAgeDays);
el('maximum-article-age').addEventListener('change',e=>setArticleAge(e.target.value));
el('manage-preferences').addEventListener('click',()=>showView('preferences'));
el('close-preferences').addEventListener('click',()=>showView('news'));
el('manage-admin').addEventListener('click',()=>{if(syncConnected&&isAdmin){showView('admin');loadAdmin()}});
el('close-admin').addEventListener('click',()=>confirmAdminLeave(()=>showView('news')));
el('admin-exit').addEventListener('click',()=>confirmAdminLeave(async()=>{
 if(!syncConnected||!isAdmin)return;
 // Exiting admin mode removes the browser's saved credential, unlike Close.
 // Require explicit confirmation even if all settings are already saved.
 const dialog=el('admin-exit-dialog');
 if(typeof dialog.showModal!=='function')return; // Fail closed if safe confirmation is unavailable.
 const confirmed=await new Promise(resolve=>{
  const done=()=>resolve(dialog.returnValue==='confirm');
  dialog.addEventListener('close',done,{once:true});
  dialog.showModal();
 });
 if(!confirmed)return;
 el('sync-disconnect').click();
 adminSavedValues=null;
 showView('news');
}));
el('admin-reload').addEventListener('click',()=>confirmAdminLeave(loadAdmin));
el('admin-form').addEventListener('submit',async event=>{event.preventDefault();await saveAdmin()});
window.addEventListener('beforeunload',event=>{if(adminDirty()){event.preventDefault();event.returnValue=''}});
syncAppearance();
const rememberedSyncToken=localStorage.getItem(syncCredentialKey);
if(rememberedSyncToken)connectSync(rememberedSyncToken,false);
else restoreReadingSelection(localStorage.getItem(readingSelectionKey));
setInterval(checkSync,30000);
document.addEventListener('visibilitychange',()=>{if(!document.hidden)checkSync()});
window.addEventListener('online',()=>{if(syncDirty&&!syncConflict)scheduleServerSave();else checkSync()});

message(storageError||exclusionError||'Search by topic or publisher to discover subscriptions.');renderExclusionManager();renderManager();showView('news');el('catalog-message').textContent='Enter a topic or publisher, then click Search or press Enter.';if(!rememberedSyncToken)renderPersonalArticles();

})();`
