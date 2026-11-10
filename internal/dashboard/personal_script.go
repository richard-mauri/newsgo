package dashboard

const personalJS = `(function(){
'use strict';
const catalog=Array.from(document.querySelectorAll('[data-included-source]'),b=>({name:b.dataset.includedSource,url:b.dataset.feedUrl,category:b.dataset.category||'Other'}));
const catalogByName=new Map(catalog.map(f=>[f.name,f]));
const feedIdentity=u=>{try{const x=new URL(u);x.hash='';return 'feed:'+x.href}catch(_){return 'feed:'+u}};
const feedRecord=f=>({id:feedIdentity(f.url),kind:'feed',name:f.name,url:f.url,category:f.category||'Other'});
const recordValid=r=>r&&typeof r.id==='string'&&typeof r.name==='string'&&r.name.length>0&&r.name.length<=140&&
 (r.kind==='feed'&&typeof r.url==='string'&&/^https?:\/\/[^\s/]+/i.test(r.url)&&r.url.length<=2048);
// Publisher exclusions are independent of subscription membership.
const exclusionKey='newsgo.excludedDomains.v1';
let excludedDomains=[],articleDomains=[],exclusionError='';
const articleExclusionKey='newsgo.excludedArticleDomains.v1';
const agePreferenceKey='newsgo.maxArticleAgeDays.v1';
const rssOnlyPreferenceKey='newsgo.rssMetadataOnly.v1';
const classificationCachePreferenceKey='newsgo.useClassificationCache.v1';
const leadCountKey='newsgo.maximumLeadStories.v1';
const leadOnlyKey='newsgo.leadOnlyAllSubscriptions.v1';
function normalizedLeadCount(v){const n=Number(v);return Number.isInteger(n)&&n>=1&&n<=20?n:3}
let maximumLeadStories=normalizedLeadCount(localStorage.getItem(leadCountKey));
let leadOnlyAllSubscriptions=localStorage.getItem(leadOnlyKey)==='true';
let rssMetadataOnly=localStorage.getItem(rssOnlyPreferenceKey)==='true';
let useClassificationCache=localStorage.getItem(classificationCachePreferenceKey)!=='false';
const allowedAges=[0,1,3,7,14,30,60];
function normalizedAge(v){const n=Number(v);return allowedAges.includes(n)?n:0}
let maxArticleAgeDays=normalizedAge(localStorage.getItem(agePreferenceKey)||0);
function setArticleAge(days){maxArticleAgeDays=normalizedAge(days);localStorage.setItem(agePreferenceKey,String(maxArticleAgeDays));el('maximum-article-age').value=String(maxArticleAgeDays);scheduleAccountSave();refresh();}

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
  const note=document.createElement('p');note.className='homepage-sync-note';note.textContent='Allowing the domain later will not restore removed subscriptions.';
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
 refreshExclusionViews();renderPersonalArticles();scheduleAccountSave();refresh();
 message(allow?'Allowed '+domain:'Excluded '+domain+(affected.length?' and removed '+affected.length+' subscription(s)':''));
 return true;
}
function readArticleExclusions(){try{const v=JSON.parse(localStorage.getItem(articleExclusionKey)||'[]');if(!Array.isArray(v)||v.some(x=>typeof x!=='string'||articleDomain(x)!==x))throw Error('Invalid article-domain list');articleDomains=[...new Set(v)]}catch(e){exclusionError='Article exclusions were not changed: '+e.message}}
function saveArticleExclusions(next){try{const domains=[...new Set(next)];localStorage.setItem(articleExclusionKey,JSON.stringify(domains));articleDomains=domains;renderExclusionManager();renderPersonalArticles();scheduleAccountSave();return true}catch(e){message('Could not save article exclusions: '+e.message);return false}}
function toggleArticleDomain(value){toggleDomain(value)}
function excludedURL(value){return domainIsExcluded(value)}
function readExclusions(){try{const data=JSON.parse(localStorage.getItem(exclusionKey)||'[]');if(!Array.isArray(data)||data.some(x=>typeof x!=='string'||exclusionDomain(x)!==x))throw Error('Invalid excluded-domain list');excludedDomains=[...new Set(data)]}catch(e){exclusionError='Exclusions were not changed: '+e.message}}
function saveExclusions(list){
 try{const unique=[...new Set(list)];localStorage.setItem(exclusionKey,JSON.stringify(unique));excludedDomains=unique;refreshExclusionViews();scheduleAccountSave();return true}
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
  const current=localStorage.getItem(unifiedKey);
  if(current===null){registry=[];return}
  const records=JSON.parse(current);
  if(!Array.isArray(records)||records.some(r=>!recordValid(r)))throw Error('Subscription registry contains invalid records');
  registry=[...new Map(records.map(r=>[r.id,r])).values()];
 }catch(e){storageError='Subscription storage: '+e.message+' Existing records were left unchanged.'}
}
loadSubscriptions();
let feeds=[];
function syncDerived(){
 feeds=registry.filter(r=>r.kind==='feed').map(r=>({name:r.name,url:r.url,category:r.category||'Other'}));
}
syncDerived();
function saveSubscriptions(next){
 try{
  const serialized=JSON.stringify(next);
  localStorage.setItem(unifiedKey,serialized);
  registry=next;syncDerived();scheduleAccountSave();return true;
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
 scheduleAccountSave();
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
  let availableCount=0,hiddenCount=0,unverifiedCount=0,subscribedCount=0;
  data.feeds.forEach(f=>{
   const followed=feeds.some(x=>x.url===f.url);
   const state=f.status==='available'?'available':f.status==='unavailable'?'unavailable':'unverified';
   // Keep subscribed feeds visible even if today's directory check failed.
   if(state==='unavailable'&&!followed){hiddenCount++;return}
   if(followed)subscribedCount++;else if(state==='available')availableCount++;else unverifiedCount++;
   const row=document.createElement('div');row.className='catalog-card';
   const info=document.createElement('div'),name=document.createElement('strong'),detail=document.createElement('small');
   name.textContent=f.name;detail.textContent=f.website||new URL(f.url).hostname;
   info.append(name,detail);const domain=publisherForFeed(f);markSite(info,domain);
   const health=document.createElement('small');health.style.display='block';
   health.textContent=state==='available'?'Feed checked: available':state==='unavailable'?'Current feed check: '+(f.reason||'unavailable'):'Feed not verified: '+(f.reason||'temporary problem');
   info.append(health);
   const b=button(followed?'Unsubscribe':state==='unverified'?'Try Subscribe':'Subscribe',()=>catalogSubscriptionClick(f,b));
   b.dataset.feedUrl=f.url;b.dataset.publisherDomain=domain;b.dataset.directoryStatus=state;
   b.disabled=!followed&&excludedURL(domain);
   const exclude=excludeButton(domain);row.append(info,b);if(exclude)row.append(exclude);root.append(row);
  });
  syncFollowButtons();
  if(!root.childElementCount){status.textContent='No usable feeds found for “'+q+'”. '+hiddenCount+' unavailable directory entries hidden. Try another search or Add a website below.';return}
  status.textContent=(availableCount+unverifiedCount+subscribedCount)+' feeds displayed · '+hiddenCount+' unavailable feeds hidden. '+subscribedCount+' subscribed · '+availableCount+' available · '+unverifiedCount+' not verified.';
 }catch(err){if(seq===searchSeq)status.textContent='Search failed: '+err.message+' Previous results remain visible; try again shortly.'}
 finally{searchInFlight=false;el('catalog-submit').disabled=false}
}
function queryEdited(){
 ++searchSeq; // Discard a response if its submitted query is no longer current.
 const q=el('catalog-search').value.trim();
 el('catalog-message').textContent=q.length>=2?'Press Search or Enter to search for “'+q+'”. Previous results stay visible.':'Enter at least two characters to search for news sources.';
}
function submitCatalogSearch(e){if(e)e.preventDefault();searchCatalog()}

function syncFollowButtons(){document.querySelectorAll('[data-feed-url]').forEach(b=>{
 const following=feeds.some(f=>f.url===b.dataset.feedUrl);
 const domain=b.dataset.publisherDomain||exclusionDomain(b.dataset.feedUrl);
 const state=b.dataset.directoryStatus;
 const catalog=b.closest('#catalog-list')!==null;
 b.textContent=following?(catalog?'Unsubscribe':'Subscribed'):b.dataset.subscribeError?'Failure':state==='unavailable'?'Unavailable':state==='unverified'?'Try Subscribe':'Subscribe';
 b.disabled=following?!catalog:excludedURL(domain)||state==='unavailable';
 if(!following&&excludedURL(domain))b.title='Allow this site before subscribing';else b.removeAttribute('title');
})}
// Only the last successful, explicit Read All Articles refresh may be exported.
let scraperTraceBatch=[];
function eligibleScraperTraceURLs(articles) {
 const out=[],seen=new Set();
 for(const article of articles) {
  if(!article||typeof article.url!=='string'||excludedURL(article.url)||articleExcluded(article.url))continue;
  try {
   const u=new URL(article.url);
   if((u.protocol!=='http:'&&u.protocol!=='https:')||u.username||u.password)continue;
   u.hash='';
   if(u.href.length>2048||seen.has(u.href))continue;
   seen.add(u.href);out.push(u.href);
  }catch(_){}
 }
 return out;
}
function downloadScraperTraceBatch() {
 if(!homepageIdentity||homepageIdentity.role!=='admin'||!suiteHandoff||suiteExpiresAt<=Date.now()){
  message('Download Article URLs requires a signed-in Homepage administrator.');
  updateScraperTraceExportVisibility();
  return;
 }
 if(activeReadingSelection!=='all'||leadOnlyAllSubscriptions||!scraperTraceBatch.length) {
  message('Choose Read All Articles and finish refreshing before exporting diagnostic URLs.');
  return;
 }
 const batch={schema:'newsgo.scrapergo-batch.v1',capturedAt:new Date().toISOString(),urls:scraperTraceBatch};
 const blob=new Blob([JSON.stringify(batch,null,2)+'\n'],{type:'application/json'});
 const address=URL.createObjectURL(blob),link=document.createElement('a');
 link.href=address;link.download='newsgo-scrapergo-urls.json';link.style.display='none';
 document.body.append(link);link.click();link.remove();
 setTimeout(()=>URL.revokeObjectURL(address),1000);
 message('Downloaded '+scraperTraceBatch.length+' eligible article URLs for local ScraperGo diagnostics.');
}
function updateScraperTraceExportVisibility(){
 const action=el('download-scrapergo-urls');
 action.hidden=!(homepageIdentity&&homepageIdentity.role==='admin'&&suiteHandoff&&suiteExpiresAt>Date.now()&&activeReadingSelection==='all'&&!leadOnlyAllSubscriptions&&scraperTraceBatch.length>0);
}
function chooseAllSubscriptions(leadOnly=false){
 const modeChanged=leadOnlyAllSubscriptions!==leadOnly;
 leadOnlyAllSubscriptions=leadOnly;
 localStorage.setItem(leadOnlyKey,String(leadOnlyAllSubscriptions));
 if(modeChanged)scheduleAccountSave();
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
 const select=el('read-source');
 // Replacing <option> elements clears the current value in Safari. Restore
 // the saved reading selection rather than silently converting it to All.
 select.value=activeReadingSelection&&activeReadingSelection!=='all'&&feeds.some(f=>f.url===activeReadingSelection)?'personal:'+activeReadingSelection:'';
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
 const restored=activeReadingSelection&&activeReadingSelection!=='all'&&feeds.some(f=>f.url===activeReadingSelection)?'personal:'+activeReadingSelection:'';
 const selection=readingSelection(restored);
 el('read-source').value=restored;
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
let latestReadStats=null;
let latestRawFeeds=[];
let latestTiming=null;
const thumbnailTiming=new Map();
const thumbnailResults=new Map();
function rawFeedBytes(feed){
 const encoded=feed.body||'';
 const binary=atob(encoded),bytes=new Uint8Array(binary.length);
 for(let i=0;i<binary.length;i++)bytes[i]=binary.charCodeAt(i);
 return bytes;
}
function renderRawFeeds(){
 const list=el('raw-feed-list'),count=el('raw-feed-count');
 if(!list||!count)return;
 list.replaceChildren();count.textContent='('+latestRawFeeds.length+' feeds)';
 if(!latestRawFeeds.length){list.textContent='No feed bodies captured in the most recent read.';return;}
 for(const feed of latestRawFeeds){
  const details=document.createElement('details');details.className='raw-feed-entry';
  const summary=document.createElement('summary');summary.textContent=feed.name+' — HTTP '+(feed.httpStatus||'unavailable')+(feed.cached?' (cached)':'');
  details.append(summary);
  const link=document.createElement('a');link.className='raw-feed-source';link.href=feed.url;link.target='_blank';link.rel='noopener noreferrer';link.textContent=feed.url;details.append(link);
  const note=document.createElement('p');note.className='personal-status';note.textContent='Content-Type: '+(feed.contentType||'not supplied')+' · '+(feed.body?Math.floor(feed.body.length*3/4):0)+' bytes (approximate)';details.append(note);
  const controls=document.createElement('div');controls.className='raw-feed-controls';
  const pre=document.createElement('pre');pre.className='raw-feed-body';
  // Decode only to display: downloads retain the unmodified original response bytes.
  const bytes=rawFeedBytes(feed);
  const decoded=new TextDecoder('utf-8',{fatal:false}).decode(bytes);
  pre.textContent=decoded;
  const copy=button('Copy Raw Feed',async()=>{
   try{await navigator.clipboard.writeText(decoded);el('raw-feed-status').textContent='Copied feed from '+feed.name+'.';}
   catch(e){el('raw-feed-status').textContent='Clipboard unavailable. Select text in the feed panel to copy it.';}
  });
  const download=button('Download Raw Feed',()=>{
   const type=(feed.contentType||'').toLowerCase();
   const ext=type.includes('json')?'.json':type.includes('xml')||type.includes('rss')||type.includes('atom')?'.xml':'.txt';
   const name=feed.name.toLowerCase().replace(/[^a-z0-9]+/g,'-').replace(/^-|-$/g,'').slice(0,45)||'feed';
   const url=URL.createObjectURL(new Blob([bytes],{type:feed.contentType||'application/octet-stream'}));
   const a=document.createElement('a');a.href=url;a.download=name+'-raw-feed'+ext;document.body.append(a);a.click();a.remove();setTimeout(()=>URL.revokeObjectURL(url),1000);
  });
  copy.disabled=!bytes.length;download.disabled=!bytes.length;controls.append(copy,download);details.append(controls,pre);list.append(details);
 }
}
function updateDiagnostics(){
 const panel=el('diagnostics-summary');if(!panel)return;
 panel.replaceChildren();
 if(!latestReadStats){panel.textContent='No completed read yet.';return;}
 const relevant=new Set(allPersonalArticles.filter(a=>!excludedURL(a.url)&&!articleExcluded(a.url)).map(a=>a.url));
 let loaded=0,failed=0,pending=0,notRequested=0;
 for(const url of relevant){const status=thumbnailResults.get(url);if(status==='loaded')loaded++;else if(status==='failed')failed++;else if(status==='pending')pending++;else notRequested++;}
 const values=[['Feed entries fetched',latestReadStats.fetched],['Accepted articles',latestReadStats.accepted],['Thumbnails loaded',loaded],['Thumbnail failures',failed],['Thumbnails pending',pending],['Images not requested',notRequested],['Classification cache hits',latestReadStats.cacheHits||0],['Classification cache misses',latestReadStats.cacheMisses||0],['Classification cache bypasses',latestReadStats.cacheBypasses||0],['Caching disabled by server',latestReadStats.cacheDisabled||0],['Server cache TTL',latestReadStats.cacheTTLMinutes===0?'Disabled (0 minutes)':latestReadStats.cacheTTLMinutes+' minutes'],['Cache hit rate',((latestReadStats.cacheHits||0)+(latestReadStats.cacheMisses||0))?Math.round(100*(latestReadStats.cacheHits||0)/((latestReadStats.cacheHits||0)+(latestReadStats.cacheMisses||0)))+'%':'N/A'],['Filtered by classification or age',latestReadStats.filtered],['Unprocessed (feed cap)',latestReadStats.unprocessed]];
 for(const [name,value] of values){const item=document.createElement('div');const label=document.createElement('strong');label.textContent=name;const number=document.createElement('div');number.textContent=String(value);number.style.fontSize='1.5rem';item.append(label,number);panel.append(item)}
 const timingPanel=el('diagnostics-timing');if(timingPanel){timingPanel.replaceChildren();if(latestTiming){const rows=[['Total refresh',latestTiming.totalMs],['Feed batch requests',latestTiming.batchMs],['Sum of per-feed fetch/wait durations',latestTiming.feedMs],['Sum of per-feed classification durations',latestTiming.classifyMs],['Browser render',latestTiming.renderMs]];for(const [label,value] of rows){const p=document.createElement('p');p.textContent=label+': '+Math.round(value)+' ms';timingPanel.append(p)}const completed=[...thumbnailTiming.values()].filter(x=>Number.isFinite(x.duration));if(completed.length){const average=completed.reduce((n,x)=>n+x.duration,0)/completed.length;const p=document.createElement('p');p.textContent='Thumbnail load/error event time: average '+Math.round(average)+' ms across '+completed.length+' completed images';timingPanel.append(p)}}else timingPanel.textContent='Timing available after the next completed read.';}
 renderRawFeeds();
 el('diagnostics-note').textContent='Source: '+latestReadStats.selection+' · '+latestReadStats.feeds+' feeds · '+latestReadStats.failedFeeds+' feed errors · '+latestReadStats.mode+'. Thumbnail outcomes can change after the read completes.';
}


// This recommendation is based only on classification decisions actually made
// for the selected feed, never the All subscriptions view.
function filteredSubscriptionStats(sourceName) {
 // Match the filtered audit shown to the reader, including age and domain exclusions.
 const classifiedFiltered=filteredEntries.filter(a=>a.source===sourceName).length;
 const articles=allPersonalArticles.filter(a=>a.source===sourceName);
 const excluded=articles.filter(a=>excludedURL(a.url)||articleExcluded(a.url)).length;
 const filtered=classifiedFiltered+excluded;
 const total=classifiedFiltered+articles.length;
 return {filtered,total,offer:total>0&&filtered/total>0.9};
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
 if(!feed)return;
 const stats=filteredSubscriptionStats(feed.name);
 if(!stats.offer)return;
 const pct=Math.round(100*stats.filtered/stats.total);
 const message=document.createElement('span');
 message.textContent=pct+'% of '+stats.total+' feed entries from '+feed.name+
  ' were filtered ('+stats.filtered+' filtered). ';
 const action=button('Unsubscribe from '+feed.name,async()=>{
  // Confirmation is required; displaying a recommendation never alters a subscription.
  if(!await confirmSubscriptionRemoval(feed.name))return;
  const record=registry.find(r=>r.kind==='feed'&&r.url===feed.url);
  if(!record||!unsubscribeRecord(record.id))return;
  // Keep the reading view unselected after removal; never load all feeds implicitly.
  ++refreshSequence;
  activeReadingSelection=null;selectedSource='';
  localStorage.removeItem(readingSelectionKey);scheduleAccountSave();
  el('read-source').value='';
  allPersonalArticles=[];personalErrors=[];filteredEntries=[];notProcessedEntries=[];
  el('personal-loading').textContent='Subscription removed. Choose another subscription to read.';
  renderManager();renderReading();
 });
 action.className='subscription-unsubscribe-button';
 action.style.cssText='margin-left:0.75rem;background:#18252d;color:#f1f5f9;border:1px solid #354650;border-radius:12px;padding:0.75rem 1.1rem;font-weight:700;font-size:1rem;cursor:pointer';
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
 const groups=new Map();
 for(const a of display){const classification=a.reason||'unknown';if(!groups.has(classification))groups.set(classification,[]);groups.get(classification).push(a)}
 if(pending.length)groups.set('unprocessed',pending);
 const labelForReason=reason=>({'too-old':'Outside age limit','homepage':'Homepage / index','unreadable':'Unreadable','unavailable':'Unavailable','unknown':'Unknown','publisher exclusion':'Publisher excluded','article-domain exclusion':'Article domain excluded','unprocessed':'Unprocessed'}[reason]||reason);
 const sorted=[...groups.entries()].sort((a,b)=>a[0].localeCompare(b[0]));
 const tabs=document.createElement('div');tabs.className='filtered-tabs';tabs.setAttribute('role','tablist');tabs.setAttribute('aria-label','Filtered article classifications');
 const panels=document.createElement('div');
 const views=[];
 const activate=index=>{views.forEach((entry,i)=>{entry.tab.setAttribute('aria-selected',String(i===index));entry.tab.tabIndex=i===index?0:-1;entry.panel.hidden=i!==index;});};
 sorted.forEach(([reason,items],index)=>{
  const tab=document.createElement('button');tab.type='button';tab.className='filtered-tab';tab.id='filtered-tab-'+index;tab.setAttribute('role','tab');tab.textContent=labelForReason(reason)+' ('+(reason==='unprocessed'?items.reduce((n,x)=>n+x.count,0):items.length)+')';
  const section=document.createElement('section');section.className='filtered-classification-group';section.id='filtered-tab-panel-'+index;section.setAttribute('role','tabpanel');section.setAttribute('aria-labelledby',tab.id);
  tab.setAttribute('aria-controls',section.id);tab.addEventListener('click',()=>activate(index));
  tab.addEventListener('keydown',event=>{if(!['ArrowLeft','ArrowRight','Home','End'].includes(event.key))return;event.preventDefault();let next=event.key==='Home'?0:event.key==='End'?sorted.length-1:(index+(event.key==='ArrowRight'?1:-1)+sorted.length)%sorted.length;activate(next);views[next].tab.focus();});
  tabs.append(tab);panels.append(section);views.push({tab,panel:section});
  if(reason==='unprocessed'){for(const x of items){const p=document.createElement('p');p.textContent=x.source+': '+x.count+' additional feed entries not processed (per-feed cap)';section.append(p)}return;}
  for(const a of items){const row=document.createElement('p');row.className='personal-row';const link=document.createElement('a');link.href=a.url;link.target='_blank';link.rel='noopener noreferrer';link.textContent=a.title||a.url;const reason=document.createElement('small');reason.textContent=' '+a.source+' — '+(a.reason==='too-old'?'Older than '+maxArticleAgeDays+' days':(a.reason||'unknown'))+(a.published&&a.reason==='too-old'?' · Published '+new Date(a.published).toLocaleDateString():'');row.append(link,reason);
 if(mayAdminister()&&a.url&&a.reason!=='too-old'&&a.reason!=='publisher exclusion'&&a.reason!=='article-domain exclusion'){
  const details=document.createElement('details');
  const summary=document.createElement('summary');summary.textContent='Admin diagnostics (run fresh classification)';
  const output=document.createElement('pre');output.textContent='Open to inspect a fresh classification request.';
  details.append(summary,output);
  details.addEventListener('toggle',async()=>{
   if(!details.open)return;
   output.textContent='Checking…';
   try{
    const response=await fetch('/api/admin/classification-diagnostic',{method:'POST',headers:{...adminAuthHeader(),'Content-Type':'application/json'},body:JSON.stringify({url:a.url}),cache:'no-store'});
    if(!response.ok)throw Error('HTTP '+response.status);
    const result=await response.json();
    const lines=['Classification: '+(result.classification||'unknown')];
    if(result.reason)lines.push('Reason: '+result.reason);
    if(Number.isInteger(result.publisherHttpStatus)&&result.publisherHttpStatus>=100&&result.publisherHttpStatus<=599)lines.push('Publisher HTTP status: '+result.publisherHttpStatus);
    if(result.failure&&result.failure!=='none')lines.push('Diagnostic failure: '+result.failure);
    if(Number.isFinite(result.elapsedMs))lines.push('Classification request duration: '+result.elapsedMs+' ms');
    output.textContent=lines.join('\n');
   }catch(error){output.textContent='Diagnostic unavailable: '+error.message;}
  });
  row.append(details);
 }
 section.append(row)}
 });
 if(views.length){list.append(tabs,panels);activate(0)}
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
  const featured=count++<maximumLeadStories;
  if(!featured&&leadOnlyAllSubscriptions&&activeReadingSelection==='all')return;
  const row=document.createElement('article');row.className=featured?'feature':'item';
  const body=document.createElement('div');body.className='item-body';
  const meta=document.createElement('div');meta.className='source';meta.textContent=a.source+' · '+a.category;
  const title=document.createElement('a');title.className='story-title';title.href=a.open;title.target='newsgo-scrapergo';title.textContent=a.title;
  const sub=document.createElement('div');sub.className='story-meta';sub.textContent=a.date+'  ';
  if(a.dateWarning){const warning=document.createElement('span');warning.className='date-warning';warning.textContent=' ⚠ '+a.dateWarning;warning.title='Original feed date shown. Article age could not be determined.';sub.append(warning);}
  const original=document.createElement('a');original.href=a.url;original.target='_blank';original.rel='noopener noreferrer';original.textContent=featured?'Original publisher ↗':'Original ↗';original.className='original';
  sub.append(original);body.append(meta,title,sub)
  if(a.image&&/^https:\/\/[^\s/]+/i.test(a.image)){
   const img=document.createElement('img');img.className='article-thumb';img.alt='';img.loading='eager';img.referrerPolicy='no-referrer';img.src=a.image;
   const imageStart=performance.now();
   if(!thumbnailResults.has(a.url))thumbnailResults.set(a.url,'pending');
   img.addEventListener('load',()=>{thumbnailResults.set(a.url,'loaded');thumbnailTiming.set(a.url,{duration:performance.now()-imageStart});updateDiagnostics()});
   img.addEventListener('error',()=>{thumbnailResults.set(a.url,'failed');thumbnailTiming.set(a.url,{duration:performance.now()-imageStart});updateDiagnostics();img.remove()});
   if(featured)body.prepend(img);else row.append(img);
  }
  if(featured){row.append(body);lead.append(row)}else{row.prepend(body);list.append(row)}
 });
 renderFilteredArticles();
 updateDiagnostics();
 const hasArticles=count>0;
 el('personal-lead-section').hidden=!hasArticles;
 el('personal-latest-section').hidden=!hasArticles||(leadOnlyAllSubscriptions&&activeReadingSelection==='all');
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
 const allActions=document.createElement('div');allActions.className='personal-form';
  allActions.append(button('Read All Articles',()=>chooseAllSubscriptions(false)),button('Read Lead Articles',()=>chooseAllSubscriptions(true)));
  allRow.append(allInfo,allActions);root.append(allRow);
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
 el('news-content').classList.toggle('personal-hidden',view!=='news');
 el('admin-panel').classList.toggle('personal-hidden',view!=='admin');
 el('diagnostics-panel').classList.toggle('personal-hidden',view!=='diagnostics');
 el('preferences-panel').classList.toggle('personal-hidden',view!=='preferences');
 [['manage-personal','discover'],['view-subscriptions','subscriptions'],['manage-preferences','preferences'],['manage-admin','admin'],['manage-diagnostics','diagnostics']].forEach(([id,value])=>el(id).setAttribute('aria-pressed',String(view===value)));
}
let refreshSequence=0;
async function refresh(){
 const refreshStart=performance.now();let batchDuration=0;
 const sequence=++refreshSequence;
 scraperTraceBatch=[];
 updateScraperTraceExportVisibility();
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
   const batchStart=performance.now();
   const res=await fetch('/api/personal',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({feeds:batch,maxArticleAgeDays,rssMetadataOnly,useClassificationCache})});
   if(sequence!==refreshSequence)return;
   if(!res.ok)throw new Error('HTTP '+res.status);
   data.push(...await res.json());batchDuration+=performance.now()-batchStart;
   if(sequence!==refreshSequence)return;
   completed+=batch.length;
   if(completed<snapshot.length)loading.textContent='Updated '+completed+' of '+snapshot.length+' web feeds…'+context;
  }
  if(sequence!==refreshSequence)return;
  allPersonalArticles=[];personalErrors=[];filteredEntries=[];notProcessedEntries=[];thumbnailResults.clear();thumbnailTiming.clear();
  latestRawFeeds=data.filter(g=>g.rawFeed).map(g=>({name:g.source.name,...g.rawFeed}));
  latestReadStats={mode:rssMetadataOnly?'RSS metadata + ScraperGo verified':'ScraperGo verified',fetched:data.reduce((n,g)=>n+(g.feedEntries||0),0),accepted:data.reduce((n,g)=>n+(g.articles||[]).length,0),filtered:data.reduce((n,g)=>n+(g.filtered||[]).length,0),unprocessed:data.reduce((n,g)=>n+(g.notProcessed||0),0),failedFeeds:data.filter(g=>g.error).length,cacheHits:data.reduce((n,g)=>n+(g.classificationCacheHits||0),0),cacheMisses:data.reduce((n,g)=>n+(g.classificationCacheMisses||0),0),cacheBypasses:data.reduce((n,g)=>n+(g.classificationCacheBypasses||0),0),cacheDisabled:data.reduce((n,g)=>n+(g.classificationCacheDisabled||0),0),cacheTTLMinutes:data.find(g=>Number.isInteger(g.classificationCacheTTLMinutes))?.classificationCacheTTLMinutes??0,feeds:snapshot.length,selection:activeReadingSelection==='all'?'All subscriptions':(snapshot[0]?.name||'One subscription')};
  data.forEach(group=>{if(group.filtered)filteredEntries.push(...group.filtered);if(group.notProcessed)notProcessedEntries.push({source:group.source.name,count:group.notProcessed})});
  data.forEach(group=>{if(group.error)personalErrors.push(group.source.name+': '+group.error);else allPersonalArticles.push(...(group.articles||[]))});
   if(activeReadingSelection==='all'&&!leadOnlyAllSubscriptions) scraperTraceBatch=eligibleScraperTraceURLs(allPersonalArticles);
   updateScraperTraceExportVisibility();
  loading.classList.remove('is-refreshing');
  const renderStart=performance.now();renderPersonalArticles();
  latestTiming={totalMs:performance.now()-refreshStart,batchMs:batchDuration,feedMs:data.reduce((n,g)=>n+(g.feedDurationMs||0),0),classifyMs:data.reduce((n,g)=>n+(g.classificationDurationMs||0),0),renderMs:performance.now()-renderStart};updateDiagnostics();
  const succeeded=data.filter(group=>!group.error).length;
  loading.textContent='Updated '+succeeded+' of '+snapshot.length+' web feeds'+(personalErrors.length?' · '+personalErrors.length+' failed':'')+context+' · '+allPersonalArticles.filter(a=>!excludedURL(a.url)&&!articleExcluded(a.url)).length+' visible headlines';
 }catch(e){
  if(sequence!==refreshSequence)return;
  loading.classList.remove('is-refreshing');
  loading.textContent='Could not refresh feeds: '+e.message;
 }
}

// Only show safe, identifying feed URL components in diagnostics. Directory
// URLs are untrusted and may include access tokens or other private query data.
function safeSubscriptionFeedURL(raw){
 try{
  const u=new URL(raw);
  if(u.protocol!=='https:'&&u.protocol!=='http:')return '(invalid feed URL)';
  u.username='';u.password='';u.hash='';
  for(const [key,value] of u.searchParams){
   if(!/^(rss|embedThumb|summary|format|type|output)$/i.test(key)||!/^(y|n|yes|no|true|false|0|1|rss|atom|xml|json)$/i.test(value)){
    u.searchParams.set(key,'[redacted]');
   }
  }
  return u.toString();
 }catch(e){return '(invalid feed URL)'}
}
function showSubscriptionFailure(trigger,f){
 const previous=document.activeElement;
 const shade=document.createElement('div');shade.className='sync-choice-backdrop';
 const dialog=document.createElement('section');dialog.className='sync-choice-dialog';
 dialog.setAttribute('role','dialog');dialog.setAttribute('aria-modal','true');
 dialog.setAttribute('aria-labelledby','subscription-failure-title');
 const heading=document.createElement('h3');heading.id='subscription-failure-title';heading.textContent='Subscription failed: '+f.name;
 const endpoint=document.createElement('p');endpoint.textContent='Attempted RSS/Atom/JSON feed URL: '+safeSubscriptionFeedURL(f.url);
 endpoint.style.overflowWrap='anywhere';
 const detail=document.createElement('p');detail.textContent='Validation failure: '+(trigger.dataset.subscribeError||'The feed could not be added.');
 const note=document.createElement('p');note.textContent='Sensitive URL query values are redacted. The address above is the attempted feed endpoint, not the publisher homepage.';note.className='personal-status';
 const actions=document.createElement('div');actions.className='sync-choice-actions';
 function finish(){document.removeEventListener('keydown',onKey);shade.remove();if(previous&&previous.isConnected)previous.focus()}
 function onKey(e){if(e.key==='Escape'){e.preventDefault();finish()}}
 const close=button('Close',finish);
 const retry=button('Retry',()=>{finish();delete trigger.dataset.subscribeError;add(f,trigger)});
 actions.append(close,retry);dialog.append(heading,endpoint,detail,note,actions);shade.append(dialog);document.body.append(shade);
 document.addEventListener('keydown',onKey);close.focus();
}
// Search buttons use the same persisted feed records as My Subscriptions.
async function catalogSubscriptionClick(candidate,trigger){
 const record=registry.find(r=>r.kind==='feed'&&r.url===candidate.url);
 if(!record){subscribeClick({name:candidate.name,category:'Other',url:candidate.url},trigger);return}
 if(!await confirmSubscriptionRemoval(record.name))return;
 if(!unsubscribeRecord(record.id))return;
 delete trigger.dataset.subscribeError;
 if(selectedSource===candidate.url||activeReadingSelection===candidate.url){
  ++refreshSequence;
  activeReadingSelection=null;selectedSource='';
  localStorage.removeItem(readingSelectionKey);scheduleAccountSave();
  el('read-source').value='';
  allPersonalArticles=[];personalErrors=[];filteredEntries=[];notProcessedEntries=[];
  el('personal-loading').textContent='Subscription removed. Choose another subscription to read.';
 }
 renderManager();renderReading();
 // A now-unsubscribed, verified-broken candidate no longer belongs in search.
 if(trigger.dataset.directoryStatus==='unavailable')trigger.closest('.catalog-card')?.remove();
 if(activeReadingSelection!==null)refresh();
 message('Unsubscribed from '+record.name+'.');
}
function subscribeClick(f,trigger){
 if(trigger&&trigger.dataset.subscribeError){showSubscriptionFailure(trigger,f);return}
 add(f,trigger);
}
function add(f,trigger){if(!valid(f)){message('Enter a source name and valid HTTP(S) feed URL.');return}if(excludedURL(trigger?.dataset.publisherDomain||publisherForFeed(f))){message('Allow this site before subscribing.');syncFollowButtons();return}if(feeds.some(x=>x.url===f.url)){message('That feed is already in your personal list.');return}
 if(trigger){delete trigger.dataset.subscribeError;trigger.disabled=true;trigger.textContent='Checking…'}
 message('Validating feed…');
 fetch('/api/personal',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({feeds:[f],maxArticleAgeDays,rssMetadataOnly,useClassificationCache})}).then(async res=>{if(!res.ok){const detail=await res.text();throw Error('NewsGo validation returned HTTP '+res.status+(detail?': '+detail.slice(0,160):''))}const data=await res.json();if(data[0].error)throw Error(/HTTP 5\d\d/.test(data[0].error)?'The publisher feed server failed ('+data[0].error+'). Try another feed result or use Add a website.':data[0].error);if(excludedURL(trigger?.dataset.publisherDomain||publisherForFeed(f)))throw Error('This site was excluded while the subscription was being checked. Allow the site before subscribing.');if(!subscribeRecord(feedRecord(f)))return;renderManager();message('Added '+f.name+' — '+(data[0].articles||[]).length+' articles found.');refresh()}).catch(err=>{const error='Feed not added: '+err.message;message(error);syncFollowButtons();if(trigger){trigger.disabled=false;trigger.dataset.subscribeError=error;trigger.textContent='Failure';trigger.title='Click to see why the subscription failed';trigger.setAttribute('aria-label','Subscription failed for '+f.name+'. Click for details.')}})}
async function discoverFeeds(){const u=el('discover-url').value.trim(),out=el('discover-results');out.replaceChildren();message('Searching for RSS or Atom feeds…');
 try{const r=await fetch('/api/discover',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({url:u})});const data=await r.json();if(!r.ok)throw Error(data.error||'HTTP '+r.status);if(!data.feeds||!data.feeds.length)throw Error('No RSS or Atom feeds found. Try Advanced with a direct feed URL.');
 data.feeds.forEach(f=>{const row=document.createElement('div');row.className='personal-row';const info=document.createElement('div');const title=document.createElement('strong');title.textContent=f.name;const small=document.createElement('small');small.textContent=(data.counts&&data.counts[f.url]||0)+' articles · '+(f.category||'News')+' · '+(new URL(f.url)).hostname;const details=document.createElement('details');const summary=document.createElement('summary');summary.textContent='Technical feed address';const address=document.createElement('small');address.textContent=f.url;details.append(summary,address);info.append(title,document.createElement('br'),small,details);const domain=publisherForFeed(f);markSite(info,domain);const subscribeButton=button('Subscribe',()=>subscribeClick(f,subscribeButton));row.append(info,subscribeButton);const exclusion=excludeButton(domain);if(exclusion)row.append(exclusion);out.append(row)});message(data.feeds.length+' feed(s) found. Select Subscribe to add one.');
 }catch(err){message('Discovery: '+err.message)}}
el('discover-button').addEventListener('click',discoverFeeds);
el('discover-form').addEventListener('submit',e=>{e.preventDefault();discoverFeeds()});
el('catalog-search').addEventListener('input',queryEdited);
el('catalog-form').addEventListener('submit',submitCatalogSearch);
el('manage-personal').addEventListener('click',()=>showView('discover'));

function chooseIncluded(name){const f=catalogByName.get(name);if(f)chooseSource(f.url)}
document.querySelectorAll('.read-included').forEach(b=>b.addEventListener('click',()=>{
 const f=catalogByName.get(b.dataset.includedSource);
 if(f)subscribeClick(f,b);
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
   const entries=records.map(r=>r&&r.kind?r:feedRecord(r));
 if(entries.some(r=>!recordValid(r)))throw Error('Invalid subscription entry');
 const unique=[...new Map(entries.map(r=>[r.id,r])).values()];
 if(!saveSubscriptions(unique))return;
   if(importedExclusions!==null&&!saveExclusions(importedExclusions))return;
   if(importedArticles!==null&&!saveArticleExclusions(importedArticles))return;
   renderManager();refresh();message('Subscriptions imported'+(importedExclusions!==null?' with exclusions':' (existing exclusions preserved)')+'.')}catch(err){message('Import failed: '+err.message)}e.target.value=''});
el('exclude-add').addEventListener('click',()=>{excludeDomain(el('exclude-domain').value);el('exclude-domain').value=''});
el('exclude-domain').addEventListener('keydown',e=>{if(e.key==='Enter'){e.preventDefault();el('exclude-add').click()}});

// Homepage is the sole NewsGo identity and profile storage provider.
const homepageOrigin='https://home.themauris.org';
const homepageSessionStatus=document.getElementById('homepage-session-status');
let homepageIdentity=null;
let suiteHandoff='';
const suiteHandoffKey='newsgo.homepage.handoff.v1';
// The local-mode fragment suppresses recovery for this visit only. A later
// direct visit must recheck Homepage, which may have gained a valid login.
let explicitBrowserOnly=false;
function readSuiteHandoff(){
 const fragment=new URLSearchParams(location.hash.slice(1));
 const delivered=fragment.get('suite_token');
 const browserOnly=fragment.get('suite_mode')==='local';
 if(delivered!==null || browserOnly){
  history.replaceState(null,'',location.pathname+location.search);
  if(browserOnly){
   // Homepage explicitly said there is no live account. Suppress redirect
   // during this page visit, preserving browser-local subscriptions.
   explicitBrowserOnly=true;
   try{sessionStorage.removeItem(suiteHandoffKey)}catch(_){}
  }else if(/^[a-f0-9]{64}$/.test(delivered)){
   try{sessionStorage.setItem(suiteHandoffKey,delivered)}catch(_){}
  }else {
   try{sessionStorage.removeItem(suiteHandoffKey)}catch(_){}
  }
 }
 try{suiteHandoff=sessionStorage.getItem(suiteHandoffKey)||''}catch(_){suiteHandoff=''}
 if(!/^[a-f0-9]{64}$/.test(suiteHandoff))suiteHandoff='';
}
readSuiteHandoff();
const homepageAccountHeaders=()=>suiteHandoff?{Authorization:'Bearer '+suiteHandoff}:{};
function homepageAdminActive(){return !!(homepageIdentity&&homepageIdentity.role==='admin'&&suiteHandoff&&suiteExpiresAt>Date.now())}
function mayAdminister(){return homepageAdminActive()}
function adminAuthHeader(){return homepageAdminActive()?{'X-NewsGo-Homepage-Authorization':'Bearer '+suiteHandoff}:{}}
let suiteExpiresAt=0;
let suiteAccountLine=null;
// Homepage alone owns countdown and renewal. NewsGo holds only a short-lived
// parent-bound handoff and asks Homepage whether it remains authorized.
function showSuiteSessionStatus(){
 if(!suiteAccountLine)return;
 if(homepageIdentity){
  suiteAccountLine.textContent='Signed in as '+homepageIdentity.user+' (Homepage)'+(accountNotice?' · '+accountNotice:'');
 }else{
  suiteAccountLine.textContent=suiteExpiresAt>0&&suiteExpiresAt<=Date.now()?'Homepage session expired':'Not signed in';
 }
}
function initSuiteAccountHeader(){
 if(suiteAccountLine||!homepageSessionStatus)return;
 const container=document.createElement('div');
 container.style.cssText='display:flex;align-items:center;flex-wrap:wrap;gap:10px;margin:5px 0 8px;color:var(--accent);font-size:.9rem;font-weight:650';
 suiteAccountLine=document.createElement('span');
 suiteAccountLine.setAttribute('aria-live','off');
 container.append(suiteAccountLine);
 homepageSessionStatus.parentNode.insertBefore(container,homepageSessionStatus);
 homepageSessionStatus.hidden=true;
 showSuiteSessionStatus();
 // A parent may have been renewed in Homepage while NewsGo was open.
 // Revalidate the child bearer at its cached deadline before discarding it.
 let deadlineCheckPending=false;
 setInterval(()=>{
  if(!suiteHandoff||!homepageIdentity||suiteExpiresAt>Date.now()||deadlineCheckPending)return;
  deadlineCheckPending=true;
  void checkHomepageSession().finally(()=>{deadlineCheckPending=false;});
 },1000);
}
function accountSessionState(){
 updateScraperTraceExportVisibility();
 updateAccountControls();
 if(!homepageIdentity){
  clearTimeout(accountTimer);accountTimer=null;
  accountReady=false;accountRevision=null;accountOwner=null;accountConflict=false;
  accountStatus('');
 }
}
async function checkHomepageSession(){
 if(!homepageSessionStatus)return;
 initSuiteAccountHeader();
 const show=t=>{homepageSessionStatus.textContent=t;homepageSessionStatus.hidden=true;homepageSessionStatus.classList.toggle('signed-in',!!homepageIdentity);showSuiteSessionStatus()};
 if(location.protocol!=='https:' || location.hostname!=='newsgo.themauris.org'){
  homepageIdentity=null;accountSessionState();show('Not signed in · Open NewsGo from Homepage');return;
 }
 if(!suiteHandoff && !explicitBrowserOnly){
  // Only Homepage can read its persistent host-only session cookie. A single
  // top-level navigation lets it issue a new parent-bound handoff (or return
  // #suite_mode=local if the login has expired/revoked). No credentials are
  // exposed to NewsGo and no browser storage grants independent sign-in.
  location.replace(homepageOrigin+'/launch/newsgo');
  return;
 }
 if(!suiteHandoff){
  homepageIdentity=null;accountSessionState();show('Not signed in · Browser settings only');return;
 }
 try{
  const verifiedToken=suiteHandoff;
  const response=await fetch(homepageOrigin+'/api/session/profile',{credentials:'omit',cache:'no-store',headers:{Authorization:'Bearer '+verifiedToken}});
  // Ignore stale verification completing after a successful bearer rotation.
  if(verifiedToken!==suiteHandoff)return;
  if(response.status===401){
   suiteExpiresAt=0;
   suiteHandoff='';try{sessionStorage.removeItem(suiteHandoffKey)}catch(_){}
   homepageIdentity=null;accountSessionState();show('Account access expired · Open NewsGo from Homepage');return;
  }
  if(!response.ok)throw Error('HTTP '+response.status);
  const profile=await response.json();
  suiteExpiresAt=Date.parse(profile.expiresAt)||0;
  if(!Number.isFinite(suiteExpiresAt)||suiteExpiresAt<=Date.now())throw Error('Expired profile');
  if(typeof profile.user!=='string'||!profile.user)throw Error('Invalid profile');
  if(homepageIdentity&&homepageIdentity.user!==profile.user){homepageIdentity=null;accountSessionState()}
  homepageIdentity={user:profile.user,role:profile.role==='admin'?'admin':'user'};
  accountSessionState();show('Signed in as '+homepageIdentity.user+' (Homepage)');showSuiteSessionStatus();
  void initializeHomepageAccount(profile.user,verifiedToken);
 }catch(e){
  homepageIdentity=null;accountSessionState();show('Account connection unavailable · Browser settings only');showSuiteSessionStatus();
 }
}
// Profile storage is user-scoped by Homepage; no user identifier is supplied.
async function homepageReaderRecord(){
 const response=await fetch(homepageOrigin+'/api/storage/newsgo',{credentials:'omit',cache:'no-store',headers:homepageAccountHeaders()});
 if(response.status===404)return null;
 if(!response.ok)throw Error('Homepage storage unavailable (HTTP '+response.status+')');
 return response.json();
}
// Homepage-verified NewsGo preferences hydrate automatically; local browser
// settings remain available when signed out. Saves use revisioned writes.
let accountRevision=null,accountReady=false,accountPending=false,accountDirty=false,accountConflict=false,accountTimer=null,accountOwner=null;
let accountNotice='';
const accountStatus=message=>{accountNotice=message||'';if(suiteAccountLine){suiteAccountLine.title=accountNotice;showSuiteSessionStatus()}};
function accountSnapshot(){return {subscriptions:registry,excludedDomains,excludedArticleDomains:articleDomains,maxArticleAgeDays,maximumLeadStories,leadOnlyAllSubscriptions,useClassificationCache,lastReadingSelection:validReadingSelection(activeReadingSelection)||''}}
function validAccountData(d){return d&&Array.isArray(d.subscriptions)&&Array.isArray(d.excludedDomains)&&Array.isArray(d.excludedArticleDomains)&&d.subscriptions.every(recordValid)&&d.excludedDomains.every(x=>typeof x==='string'&&exclusionDomain(x)===x)&&d.excludedArticleDomains.every(x=>typeof x==='string'&&articleDomain(x)===x)}
async function accountSave(){
 await checkHomepageSession();
 if(!homepageIdentity||!accountReady||accountPending||accountConflict||!accountDirty||accountOwner!==homepageIdentity.user)return;
 accountPending=true;accountDirty=false;
 const owner=accountOwner,token=suiteHandoff,revision=accountRevision;
 try{
  const data=accountSnapshot();
  if(!validAccountData(data))throw Error('Local profile failed validation');
  const response=await fetch(homepageOrigin+'/api/storage/newsgo',{method:'PUT',credentials:'omit',cache:'no-store',headers:{'Content-Type':'application/json',...homepageAccountHeaders()},body:JSON.stringify({revision,data})});
  // An older in-flight write must never update a newly signed-in user's
  // revision or account status.
  if(!homepageIdentity||homepageIdentity.user!==owner||accountOwner!==owner||suiteHandoff!==token)return;
  if(response.status===409){accountConflict=true;throw Error('Account settings changed elsewhere. Local changes were not saved; reload to use the latest account record.')}
  if(!response.ok)throw Error('Account save failed (HTTP '+response.status+')');
  const result=await response.json();
  if(!homepageIdentity||homepageIdentity.user!==owner||accountOwner!==owner||suiteHandoff!==token)return;
  accountRevision=result.revision;
  accountStatus('');
 }catch(e){if(accountOwner===owner){accountDirty=true;accountStatus(e.message)}}
 finally{accountPending=false;if(accountDirty&&!accountConflict&&accountReady&&homepageIdentity&&homepageIdentity.user===accountOwner){clearTimeout(accountTimer);accountTimer=setTimeout(accountSave,30000)}}
}
function scheduleAccountSave(){
 if(!homepageIdentity||!accountReady||accountConflict)return;
 accountDirty=true;clearTimeout(accountTimer);accountTimer=setTimeout(accountSave,650);
}
// Homepage account profile is loaded automatically after verified sign-in.
// Never treat a failed GET as an empty account; do not overwrite on 409.
let accountLoadingUser=null;
async function initializeHomepageAccount(owner,token){
 if(accountReady&&accountOwner===owner)return;
 if(accountLoadingUser===owner)return;
 accountLoadingUser=owner;
 try{
  const record=await homepageReaderRecord();
  if(!homepageIdentity||homepageIdentity.user!==owner||suiteHandoff!==token)return;
  if(record&&(!validAccountData(record.data)||!Number.isSafeInteger(record.revision)||record.revision<0))
   throw Error('Invalid saved account settings; local preferences retained');
  clearTimeout(accountTimer);accountTimer=null;
  let prior='';
  try{prior=localStorage.getItem('newsgo.account.last-owner.v1')||''}catch(_){}
  if(record){
   applyAccountPreferences(record.data);
   accountRevision=record.revision;
  }else{
   // A new account must not inherit a different verified user's local data.
   if(prior&&prior!==owner){
    applyAccountPreferences({subscriptions:[],excludedDomains:[],excludedArticleDomains:[],maxArticleAgeDays:0,maximumLeadStories:3,leadOnlyAllSubscriptions:false,useClassificationCache:true,lastReadingSelection:''});
   }
   accountRevision=0;
  }
  accountOwner=owner;accountReady=true;accountConflict=false;accountDirty=false;
  try{localStorage.setItem('newsgo.account.last-owner.v1',owner)}catch(_){}
  accountStatus('');
  if(!record)scheduleAccountSave();
 }catch(e){
  accountReady=false;
  accountStatus('Account preferences could not load: '+e.message);
 }finally{if(accountLoadingUser===owner)accountLoadingUser=null}
}
let adminRevision=0;
function updateAccountControls(){
 el('manage-admin').hidden=!homepageAdminActive();
 el('manage-diagnostics').hidden=false;
 if(!homepageAdminActive()&&!el('admin-panel').classList.contains('personal-hidden'))showView('news');
}
function applyAccountPreferences(data){
 if(!Array.isArray(data.subscriptions)||!Array.isArray(data.excludedDomains)||!Array.isArray(data.excludedArticleDomains||[]))throw Error('Invalid server profile');
 if(data.subscriptions.some(r=>!recordValid(r))||data.excludedDomains.some(d=>typeof d!=='string'||exclusionDomain(d)!==d)||!Array.isArray(data.excludedArticleDomains||[])||(data.excludedArticleDomains||[]).some(d=>typeof d!=='string'||articleDomain(d)!==d))throw Error('Server profile contains invalid entries');
 const entries=[...new Map(data.subscriptions.map(r=>[r.id,r])).values()];
 localStorage.setItem(unifiedKey,JSON.stringify(entries));localStorage.setItem(exclusionKey,JSON.stringify(data.excludedDomains));localStorage.setItem(articleExclusionKey,JSON.stringify(data.excludedArticleDomains||[]));
 registry=entries;syncDerived();excludedDomains=[...new Set(data.excludedDomains)];articleDomains=[...new Set(data.excludedArticleDomains||[])];maxArticleAgeDays=normalizedAge(data.maxArticleAgeDays||0);localStorage.setItem(agePreferenceKey,String(maxArticleAgeDays));el('maximum-article-age').value=String(maxArticleAgeDays);
  maximumLeadStories=normalizedLeadCount(data.maximumLeadStories||3);leadOnlyAllSubscriptions=data.leadOnlyAllSubscriptions===true;
  localStorage.setItem(leadCountKey,String(maximumLeadStories));localStorage.setItem(leadOnlyKey,String(leadOnlyAllSubscriptions));
  el('maximum-lead-stories').value=String(maximumLeadStories);
   useClassificationCache=data.useClassificationCache!==false;localStorage.setItem(classificationCachePreferenceKey,String(useClassificationCache));el('use-classification-cache').checked=useClassificationCache;
 const restored=validReadingSelection(data.lastReadingSelection);if(restored===null)localStorage.removeItem(readingSelectionKey);else localStorage.setItem(readingSelectionKey,restored);
 restoreReadingSelection(restored);

 refreshExclusionViews();renderExclusionManager();renderManager();refresh();
}
const articleAdminFields={maxArticlesPerFeed:'admin-articles',classificationTimeoutSeconds:'admin-classification-timeout',classificationConcurrency:'admin-classification-parallel',feedCacheMinutes:'admin-feed-cache',feedSearchTimeoutSeconds:'admin-feed-search-timeout',classificationCacheMinutes:'admin-classification-cache'};
// Add classification-cache TTL control alongside feed-cache settings.
const feedCacheInput=document.getElementById('admin-feed-cache');
if(feedCacheInput){
 const label=document.createElement('label');
 label.setAttribute('for','admin-classification-cache');
 label.textContent='Successful classification cache (minutes, 0–1440)';
 label.style.display='block';label.style.marginTop='1rem';
 const input=document.createElement('input');
 input.type='number';input.min='0';input.max='1440';input.step='1';input.required=true;
 input.id='admin-classification-cache';input.value='1440';
 input.style.display='block';input.style.width='100%';input.style.maxWidth='22rem';
 label.append(input);
 const parent=feedCacheInput.closest('label')||feedCacheInput.parentElement;
 parent.insertAdjacentElement('afterend',label);
}
let currentImageSettings=null;
let adminSavedValues=null;
const adminStatus=message=>{el('admin-status').textContent=message};
const adminValues=()=>Object.fromEntries(Object.entries(articleAdminFields).map(([key,id])=>[key,el(id).value]));
const adminDirty=()=>adminSavedValues!==null&&Object.keys(articleAdminFields).some(key=>adminValues()[key]!==adminSavedValues[key]);
async function adminRequest(method,payload){
 const response=await fetch('/api/admin/settings',{method,headers:{...adminAuthHeader(),...(payload?{'Content-Type':'application/json'}:{})},...(payload?{body:JSON.stringify(payload)}:{}),cache:'no-store'});
 if(!response.ok)throw Error('HTTP '+response.status+' — '+(await response.text()).trim().slice(0,140));return response.json();
}
async function loadAdmin(){
 if(!mayAdminister())return false;
 try{const data=await adminRequest('GET');adminRevision=data.revision;
  currentImageSettings=data.image;
  for(const [key,id] of Object.entries(articleAdminFields)){const field=el(id);if(field)field.value=(data.articles||{})[key]??'';}
  adminSavedValues=adminValues();
  adminStatus('Server settings loaded (revision '+adminRevision+').');
  return true;
 }catch(e){adminStatus('Could not load settings: '+e.message);return false}
}
async function saveAdmin(){
 if(!mayAdminister()||!currentImageSettings)return false;
 const form=el('admin-form');
 if(!form.reportValidity())return false;
 const articles=Object.fromEntries(Object.entries(adminValues()).map(([key,value])=>[key,Number(value)]));
 try{const result=await adminRequest('PUT',{revision:adminRevision,image:currentImageSettings,articles});
  adminRevision=result.revision;currentImageSettings=result.image;adminSavedValues=adminValues();
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
el('maximum-lead-stories').value=String(maximumLeadStories);
el('maximum-lead-stories').addEventListener('change',e=>{maximumLeadStories=normalizedLeadCount(e.target.value);localStorage.setItem(leadCountKey,String(maximumLeadStories));scheduleAccountSave();renderPersonalArticles()});
el('use-classification-cache').checked=useClassificationCache;
el('use-classification-cache').addEventListener('change',event=>{useClassificationCache=event.target.checked;localStorage.setItem(classificationCachePreferenceKey,String(useClassificationCache));scheduleAccountSave();});
el('rss-metadata-only').checked=rssMetadataOnly;
el('rss-metadata-only').addEventListener('change',event=>{rssMetadataOnly=event.target.checked;localStorage.setItem(rssOnlyPreferenceKey,String(rssMetadataOnly));});
el('manage-preferences').addEventListener('click',()=>showView('preferences'));
el('close-preferences').addEventListener('click',()=>showView('news'));
el('manage-admin').addEventListener('click',()=>{if(mayAdminister()){showView('admin');loadAdmin()}});
el('manage-diagnostics').addEventListener('click',()=>{showView('diagnostics');updateDiagnostics()});
el('close-diagnostics').addEventListener('click',()=>showView('news'));
el('close-admin').addEventListener('click',()=>confirmAdminLeave(()=>showView('news')));
el('admin-reload').addEventListener('click',()=>confirmAdminLeave(loadAdmin));
el('admin-form').addEventListener('submit',async event=>{event.preventDefault();await saveAdmin()});
window.addEventListener('beforeunload',event=>{if(adminDirty()){event.preventDefault();event.returnValue=''}});
el('download-scrapergo-urls').addEventListener('click',downloadScraperTraceBatch);
updateAccountControls();
restoreReadingSelection(localStorage.getItem(readingSelectionKey));
void checkHomepageSession();
window.addEventListener('focus',()=>{void checkHomepageSession()});
document.addEventListener('visibilitychange',()=>{if(!document.hidden)void checkHomepageSession()});
message(storageError||exclusionError||'Search by topic or publisher to discover subscriptions.');renderExclusionManager();renderManager();showView('news');el('catalog-message').textContent='Enter a topic or publisher, then click Search or press Enter.';renderPersonalArticles();

})();`
