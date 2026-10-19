package dashboard

const personalJS = `(function(){
'use strict';
const key='newsgo.personalFeeds.v1';
const includedKey='newsgo.includedSubscriptions.v1';
const includedNames=Array.from(document.querySelectorAll('[data-included-source]'),b=>b.dataset.includedSource);
let included=[];try{const stored=JSON.parse(localStorage.getItem(includedKey)||'[]');if(Array.isArray(stored))included=stored.filter(n=>typeof n==='string'&&includedNames.includes(n))}catch(e){}
function saveIncluded(){try{localStorage.setItem(includedKey,JSON.stringify(included))}catch(e){message('Browser storage is unavailable.')}}
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
let feeds=[];
try{const raw=JSON.parse(localStorage.getItem(key)||'[]');if(Array.isArray(raw)) feeds=raw.filter(valid).slice(0,8)}catch(e){}
function save(){try{localStorage.setItem(key,JSON.stringify(feeds));return true}catch(e){message('Browser storage is unavailable; changes will not survive a restart.');return false}}
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
   info.append(name,detail);
   const followed=feeds.some(x=>x.url===f.url);
   const b=button(followed?'Subscribed':'Subscribe',()=>add({name:f.name,category:'Other',url:f.url},b));b.dataset.feedUrl=f.url;b.disabled=followed;
   row.append(info,b);root.append(row);
  });
  status.textContent=data.feeds.length+' candidate sources. Each feed is validated when you select Subscribe.';
 }catch(err){if(seq===searchSeq)status.textContent='Search failed: '+err.message+' Try Add a website below.'}
}
function queueSearch(){clearTimeout(searchTimer);searchTimer=setTimeout(searchCatalog,350)}
function syncFollowButtons(){document.querySelectorAll('[data-feed-url]').forEach(b=>{const following=feeds.some(f=>f.url===b.dataset.feedUrl);b.textContent=following?'Subscribed':'Subscribe';b.disabled=following})}
function chooseSource(url){
 selectedSource=url;el('read-source').value='personal:'+url;renderReading();showView('news');
 el('reading-title').scrollIntoView({block:'start',behavior:'smooth'});
}
function renderSourceSelector(){
 const group=el('followed-options');group.replaceChildren();
 feeds.forEach(f=>{const option=document.createElement('option');option.value='personal:'+f.url;option.textContent=f.name;group.append(option)});
 if(selectedSource&&!feeds.some(f=>f.url===selectedSource))selectedSource='';
 const select=el('read-source');if(select.value.startsWith('personal:')&&!feeds.some(f=>'personal:'+f.url===select.value))select.value='';
}
function renderReading(){
 const source=el('read-source').value;
 selectedSource=source.startsWith('personal:')?source.slice(9):'';
 const showPersonal=!source.startsWith('shared:');
 const showShared=source.startsWith('shared:');
 el('personal-feed').classList.toggle('personal-hidden',!feeds.length||!showPersonal);
 el('shared-headlines').classList.toggle('personal-hidden',!showShared);
 document.querySelectorAll('#shared-headlines article[data-source]').forEach(a=>{
  a.hidden=!!((source&&source!=='shared:'+a.dataset.source));
 });
 el('reading-title').hidden=!source;
 el('reading-title').textContent=source?(source.startsWith('shared:')?source.slice(7):feeds.find(f=>'personal:'+f.url===source)?.name||'Selected subscription'):'';
 el('personal-section-label').hidden=!!source;
 renderPersonalArticles();
}
let allPersonalArticles=[], personalErrors=[];
function renderPersonalArticles(){
 const list=el('personal-articles');list.replaceChildren();
 const sources=new Set(feeds.filter(f=>!selectedSource||f.url===selectedSource).map(f=>f.name));
 const articles=allPersonalArticles.filter(a=>sources.has(a.source));
 const seen=new Set();articles.sort((a,b)=>new Date(b.published)-new Date(a.published));
 articles.forEach(a=>{if(seen.has(a.url))return;seen.add(a.url);const row=document.createElement('article');row.className='item';const body=document.createElement('div');body.className='item-body';const meta=document.createElement('div');meta.className='source';meta.textContent=a.source+' · '+a.category;const title=document.createElement('a');title.className='story-title';title.href=a.open;title.target='newsgo-scrapergo';title.textContent=a.title;const sub=document.createElement('div');sub.className='story-meta';sub.textContent=a.date+'  ';const original=document.createElement('a');original.href=a.url;original.target='_blank';original.rel='noopener noreferrer';original.textContent='Original ↗';original.className='original';sub.append(original);body.append(meta,title,sub);row.append(body);if(a.image&&/^https:\/\/[^\s/]+/i.test(a.image)){const img=document.createElement('img');img.className='article-thumb';img.src=a.image;img.alt='';img.loading='lazy';img.referrerPolicy='no-referrer';img.addEventListener('error',()=>img.remove());row.append(img)}list.append(row)});
 const selected=feeds.find(f=>f.url===selectedSource);
 el('personal-loading').textContent=seen.size+' headlines'+(selected?' from '+selected.name:' from your followed sources')+(personalErrors.length?' · '+personalErrors.join('; '):'');
}
function renderManager(){
 renderSourceSelector();
 const total=feeds.length+included.length;el('personal-count').textContent='('+total+')';el('subscription-count').textContent='('+total+')';
 el('subscriptions-empty').textContent=feeds.length+included.length?'':'No subscriptions yet. Choose Find Subscriptions to add one.';
 const root=el('personal-sources');root.replaceChildren();
 included.forEach(name=>{const row=document.createElement('div');row.className='personal-row';const info=document.createElement('strong');info.textContent=name;const actions=document.createElement('div');actions.className='subscription-actions';actions.append(button('Read Articles',()=>chooseIncluded(name)),button('Unsubscribe',()=>{included=included.filter(n=>n!==name);if(el('read-source').value==='shared:'+name)el('read-source').value='';saveIncluded();renderManager();}));row.append(info,actions);root.append(row)});
 document.querySelectorAll('.read-included').forEach(b=>{const subscribed=included.includes(b.dataset.includedSource);b.textContent=subscribed?'Subscribed':'Subscribe';b.disabled=subscribed});
 feeds.forEach((f,i)=>{const row=document.createElement('div');row.className='personal-row';const info=document.createElement('div');const name=document.createElement('strong');name.textContent=f.name;const detail=document.createElement('small');try{detail.textContent=new URL(f.url).hostname+' · '+f.category}catch(e){detail.textContent=f.category}info.append(name,detail);const actions=document.createElement('div');actions.className='subscription-actions';actions.append(button('Read Articles',()=>chooseSource(f.url)),button('Unsubscribe',()=>{feeds.splice(i,1);if(selectedSource===f.url)selectedSource='';save();renderManager();refresh();message('Unsubscribeed '+f.name+'.')}));row.append(info,actions);root.append(row)});
 syncFollowButtons();renderReading();
}
function showView(view){
 el('personal-panel').classList.toggle('personal-hidden',view!=='discover');
 el('subscriptions-panel').classList.toggle('personal-hidden',view!=='subscriptions');
 el('news-content').classList.toggle('personal-hidden',view!=='news');
 [['manage-personal','discover'],['view-subscriptions','subscriptions']].forEach(([id,value])=>el(id).setAttribute('aria-pressed',String(view===value)));
}
async function refresh(){const section=el('personal-feed');renderReading();if(!feeds.length){allPersonalArticles=[];renderPersonalArticles();return}
 el('personal-loading').textContent='Refreshing personal feeds…';
 try{const res=await fetch('/api/personal',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({feeds})});if(!res.ok)throw new Error('HTTP '+res.status);const data=await res.json();allPersonalArticles=[];personalErrors=[];data.forEach(group=>{if(group.error)personalErrors.push(group.source.name+': '+group.error);else allPersonalArticles.push(...(group.articles||[]))});renderPersonalArticles();
 }catch(e){el('personal-loading').textContent='Could not fetch personal feeds: '+e.message}}

function add(f,trigger){if(!valid(f)){message('Enter a source name and valid HTTP(S) feed URL.');return}if(feeds.some(x=>x.url===f.url)){message('That feed is already in your personal list.');return}if(feeds.length>=8){message('Maximum 8 personal feeds for this version.');return}
 if(trigger){trigger.disabled=true;trigger.textContent='Checking…'}
 message('Validating feed…');
 fetch('/api/personal',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({feeds:[f]})}).then(async res=>{if(!res.ok)throw Error('HTTP '+res.status);const data=await res.json();if(data[0].error)throw Error(data[0].error);feeds.push(f);save();renderManager();message('Added '+f.name+' — '+(data[0].articles||[]).length+' articles found.');refresh()}).catch(err=>{message('Feed not added: '+err.message);syncFollowButtons();if(trigger){trigger.disabled=false;trigger.textContent='Subscribe'}})}
async function discoverFeeds(){const u=el('discover-url').value.trim(),out=el('discover-results');out.replaceChildren();message('Searching for RSS or Atom feeds…');
 try{const r=await fetch('/api/discover',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({url:u})});const data=await r.json();if(!r.ok)throw Error(data.error||'HTTP '+r.status);if(!data.feeds||!data.feeds.length)throw Error('No RSS or Atom feeds found. Try Advanced with a direct feed URL.');
 data.feeds.forEach(f=>{const row=document.createElement('div');row.className='personal-row';const info=document.createElement('div');const title=document.createElement('strong');title.textContent=f.name;const small=document.createElement('small');small.textContent=(data.counts&&data.counts[f.url]||0)+' articles · '+(f.category||'News')+' · '+(new URL(f.url)).hostname;const details=document.createElement('details');const summary=document.createElement('summary');summary.textContent='Technical feed address';const address=document.createElement('small');address.textContent=f.url;details.append(summary,address);info.append(title,document.createElement('br'),small,details);row.append(info,button('Subscribe',()=>add(f)));out.append(row)});message(data.feeds.length+' feed(s) found. Select Subscribe to add one.');
 }catch(err){message('Discovery: '+err.message)}}
el('discover-button').addEventListener('click',discoverFeeds);
el('discover-form').addEventListener('submit',e=>{e.preventDefault();discoverFeeds()});
el('catalog-search').addEventListener('input',queueSearch);
el('manage-personal').addEventListener('click',()=>showView('discover'));

function chooseIncluded(name){el('read-source').value='shared:'+name;renderReading();showView('news');el('reading-title').scrollIntoView({block:'start',behavior:'smooth'})}
document.querySelectorAll('.read-included').forEach(b=>b.addEventListener('click',()=>{const name=b.dataset.includedSource;if(!included.includes(name)){included.push(name);saveIncluded();renderManager();message('Subscribed to '+name+'.')} }));
el('view-subscriptions').addEventListener('click',()=>showView('subscriptions'));
el('close-discovery').addEventListener('click',()=>showView('news'));
el('close-subscriptions').addEventListener('click',()=>showView('news'));
el('personal-form').addEventListener('submit',e=>{e.preventDefault();add({name:el('personal-name').value.trim(),url:el('personal-url').value.trim(),category:el('personal-category').value});});
const sug=el('personal-suggestions');suggestions.forEach(s=>sug.append(button('+ '+s.name,()=>add(s),'personal-suggestion')));
el('personal-export').addEventListener('click',()=>{const blob=new Blob([JSON.stringify(feeds,null,2)],{type:'application/json'}),u=URL.createObjectURL(blob),a=document.createElement('a');a.href=u;a.download='newsgo-personal-feeds.json';a.click();setTimeout(()=>URL.revokeObjectURL(u),1000)});
el('personal-import').addEventListener('click',()=>el('personal-file').click());
el('personal-file').addEventListener('change',async e=>{const f=e.target.files[0];if(!f)return;try{if(f.size>16384)throw Error('File too large');const data=JSON.parse(await f.text());if(!Array.isArray(data)||data.length>8||!data.every(valid))throw Error('Invalid feed list');feeds=data;save();renderManager();refresh();message('Subscriptions imported.')}catch(err){message('Import failed: '+err.message)}e.target.value=''});
message('Search by topic or publisher to discover subscriptions.');renderManager();showView('news');searchCatalog();refresh();setInterval(refresh,15*60*1000);
document.querySelectorAll('#shared-headlines img.article-thumb').forEach(img=>img.addEventListener('error',()=>img.remove()));
})();`
