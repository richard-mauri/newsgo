package dashboard

const personalJS = `(function(){
'use strict';
const key='newsgo.personalFeeds.v1';
const suggestions=[
{name:'NASA News',category:'Science',url:'https://www.nasa.gov/news-release/feed/'},
{name:'NASA Science',category:'Science',url:'https://science.nasa.gov/feed/'},
{name:'NASA JPL',category:'Science',url:'https://www.jpl.nasa.gov/feeds/news/'},
{name:'NOAA News',category:'Science',url:'https://www.noaa.gov/rss.xml'},
{name:'Ars Technica Science',category:'Science',url:'https://feeds.arstechnica.com/arstechnica/science'}
];
const el=id=>document.getElementById(id);
const valid=f=>f&&typeof f.name==='string'&&f.name.trim().length>0&&f.name.length<=70&&typeof f.category==='string'&&f.category.length<=40&&typeof f.url==='string'&&f.url.length<=2048&&/^https?:\/\/[^\s/]+/i.test(f.url);
let feeds=[];
try{const raw=JSON.parse(localStorage.getItem(key)||'[]');if(Array.isArray(raw)) feeds=raw.filter(valid).slice(0,8)}catch(e){}
function save(){try{localStorage.setItem(key,JSON.stringify(feeds));return true}catch(e){message('Browser storage is unavailable; changes will not survive a restart.');return false}}
function message(t){el('personal-message').textContent=t}
function button(label,handler,cls){const b=document.createElement('button');b.type='button';b.textContent=label;if(cls)b.className=cls;b.addEventListener('click',handler);return b}
function renderManager(){el('personal-count').textContent=feeds.length?'('+feeds.length+')':'';const root=el('personal-sources');root.replaceChildren();
 feeds.forEach((f,i)=>{const row=document.createElement('div');row.className='personal-row';const info=document.createElement('div');const name=document.createElement('strong');name.textContent=f.name+' · '+f.category;const detail=document.createElement('small');detail.textContent=f.url;info.append(name,document.createElement('br'),detail);row.append(info,button('Remove',()=>{feeds.splice(i,1);save();renderManager();refresh()}));root.append(row)})}
async function refresh(){const section=el('personal-feed'), list=el('personal-articles');list.replaceChildren();section.classList.toggle('personal-hidden',feeds.length===0);if(!feeds.length)return;
 el('personal-loading').textContent='Refreshing personal feeds…';
 try{const res=await fetch('/api/personal',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({feeds})});if(!res.ok)throw new Error('HTTP '+res.status);const data=await res.json();const articles=[];const errors=[];data.forEach(group=>{if(group.error)errors.push(group.source.name+': '+group.error);else articles.push(...(group.articles||[]))});
 const seen=new Set();articles.sort((a,b)=>new Date(b.published)-new Date(a.published));articles.forEach(a=>{if(seen.has(a.url))return;seen.add(a.url);const row=document.createElement('article');row.className='item';const body=document.createElement('div');body.className='item-body';const meta=document.createElement('div');meta.className='source';meta.textContent=a.source+' · '+a.category;const title=document.createElement('a');title.className='story-title';title.href=a.open;title.target='newsgo-scrapergo';title.textContent=a.title;const sub=document.createElement('div');sub.className='story-meta';sub.textContent=a.date+'  ';const original=document.createElement('a');original.href=a.url;original.target='_blank';original.rel='noopener noreferrer';original.textContent='Original ↗';original.className='original';sub.append(original);body.append(meta,title,sub);row.append(body);list.append(row)});
 el('personal-loading').textContent=seen.size+' personal headlines'+(errors.length?' · '+errors.join('; '):'');
 }catch(e){el('personal-loading').textContent='Could not fetch personal feeds: '+e.message}}
function add(f){if(!valid(f)){message('Enter a source name and valid HTTP(S) feed URL.');return}if(feeds.some(x=>x.url===f.url)){message('That feed is already in your personal list.');return}if(feeds.length>=8){message('Maximum 8 personal feeds for this version.');return}
 message('Validating feed…');
 fetch('/api/personal',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({feeds:[f]})}).then(async res=>{if(!res.ok)throw Error('HTTP '+res.status);const data=await res.json();if(data[0].error)throw Error(data[0].error);feeds.push(f);save();renderManager();message('Added '+f.name+' — '+(data[0].articles||[]).length+' articles found.');refresh()}).catch(err=>message('Feed not added: '+err.message))}
el('manage-personal').addEventListener('click',()=>el('personal-panel').classList.toggle('personal-hidden'));
el('personal-form').addEventListener('submit',e=>{e.preventDefault();add({name:el('personal-name').value.trim(),url:el('personal-url').value.trim(),category:el('personal-category').value});});
const sug=el('personal-suggestions');suggestions.forEach(s=>sug.append(button('+ '+s.name,()=>add(s),'personal-suggestion')));
el('personal-export').addEventListener('click',()=>{const blob=new Blob([JSON.stringify(feeds,null,2)],{type:'application/json'}),u=URL.createObjectURL(blob),a=document.createElement('a');a.href=u;a.download='newsgo-personal-feeds.json';a.click();setTimeout(()=>URL.revokeObjectURL(u),1000)});
el('personal-import').addEventListener('click',()=>el('personal-file').click());
el('personal-file').addEventListener('change',async e=>{const f=e.target.files[0];if(!f)return;try{if(f.size>16384)throw Error('File too large');const data=JSON.parse(await f.text());if(!Array.isArray(data)||data.length>8||!data.every(valid))throw Error('Invalid feed list');feeds=data;save();renderManager();refresh();message('Subscriptions imported.')}catch(err){message('Import failed: '+err.message)}e.target.value=''});
renderManager();refresh();setInterval(refresh,15*60*1000);
})();`
