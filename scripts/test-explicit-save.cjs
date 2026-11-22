'use strict';
// CI-only execution tests; no Node.js is required by the owner Mac build.
const fs=require('node:fs');
const vm=require('node:vm');
const assert=require('node:assert/strict');
const full=fs.readFileSync('internal/dashboard/personal_script.go','utf8');
const prefix='const personalJS = `';
const begin=full.indexOf(prefix),end=full.lastIndexOf('`');
assert.ok(begin>=0&&end>begin,'embedded NewsGo dashboard JS found');
const script=full.slice(begin+prefix.length,end);
new vm.Script(script,{filename:'newsgo-personal.js'});
const stageStart=script.indexOf("const readerUnsyncedKey=");
const stageEnd=script.indexOf('function normalizedLeadCount',stageStart);
const commitStart=script.indexOf('// Save is invoked only by the user-triggered Save Changes button.');
const commitEnd=script.indexOf('// Homepage account profile is loaded automatically',commitStart);
assert.ok(stageStart>0&&stageEnd>stageStart&&commitStart>0&&commitEnd>commitStart);
const functions=script.slice(stageStart,stageEnd)+'\n'+script.slice(commitStart,commitEnd);
function fixture({role='user',networkStatus=200,failStorage=false}={}){
 const writes=new Map(),calls=[],controls={
  'reader-save-changes':{disabled:true},
  'reader-save-status':{textContent:''}
 };
 const storage={
  getItem:k=>writes.has(k)?writes.get(k):null,
  setItem(k,v){if(failStorage&&k==='newsgo.maximumLeadStories.v1')throw Error('quota failure');writes.set(k,String(v))},
  removeItem:k=>writes.delete(k)
 };
 const snapshot={
  subscriptions:[],excludedDomains:[],excludedArticleDomains:[],
  maxArticleAgeDays:0,maximumLeadStories:4,leadOnlyAllSubscriptions:false,
  useClassificationCache:true,lastReadingSelection:''
 };
 const ctx={
  document:{getElementById:id=>controls[id]||null},
  localStorage:storage,
  homepageIdentity:role?{user:role==='guest'?'guest-001':'reader',role}:null,
  suiteHandoff:'a'.repeat(64),
  accountReady:true,accountPending:false,accountConflict:false,
  accountOwner:role==='guest'?'':role?'reader':null,accountRevision:1,
  accountNotice:'',homepageOrigin:'https://home.themauris.org',
  accountSnapshot:()=>snapshot,validAccountData:()=>true,
  accountStatus(v){ctx.accountNotice=v},
  recordAccountDiagnostic(){},
  fetch:async(url,options)=>{
   calls.push({url,options});
   return {ok:networkStatus===200,status:networkStatus,
    json:async()=>({revision:2})};
  },
  Map,JSON,Error,String,encodeURIComponent
 };
 vm.createContext(ctx);
 vm.runInContext(functions,ctx);
 return {ctx,writes,calls,controls,storage};
}
async function run(){
 let x=fixture();
 x.ctx.stageReaderSetting('newsgo.maximumLeadStories.v1','4');
 assert.equal(x.writes.size,0,'editing must not immediately persist locally');
 assert.equal(x.calls.length,0,'editing must not issue Homepage PUT');
 assert.equal(x.controls['reader-save-changes'].disabled,false);
 assert.match(x.controls['reader-save-status'].textContent,/Unsaved changes/);
 await x.ctx.commitReaderSettings();
 assert.equal(x.writes.get('newsgo.maximumLeadStories.v1'),'4');
 assert.equal(x.calls.length,1,'exactly one user-triggered Homepage PUT');
 assert.equal(x.calls[0].options.method,'PUT');
 assert.equal(x.controls['reader-save-status'].textContent,'Settings saved.');
 assert.equal(x.controls['reader-save-changes'].disabled,true);
 assert.equal(x.writes.has('newsgo.account.unsynced.v1.reader'),false);
 let g=fixture({role:'guest'});
 g.ctx.stageReaderSetting('newsgo.maximumLeadStories.v1','4');
 await g.ctx.commitReaderSettings();
 assert.equal(g.writes.get('newsgo.maximumLeadStories.v1'),'4');
 assert.equal(g.calls.length,0,'Guest must never save to Homepage account');
 assert.equal(g.controls['reader-save-status'].textContent,'Settings saved.');
 let failed=fixture({role:'user',networkStatus:503});
 failed.ctx.stageReaderSetting('newsgo.maximumLeadStories.v1','4');
 await failed.ctx.commitReaderSettings();
 assert.equal(failed.calls.length,1);
 assert.equal(failed.controls['reader-save-changes'].disabled,false,'failed PUT enables manual retry');
 assert.match(failed.controls['reader-save-status'].textContent,/retry/i);
 assert.ok(failed.writes.has('newsgo.account.unsynced.v1.reader'),'per-user retry payload persisted');
 failed.ctx.fetch=async()=>({ok:true,status:200,json:async()=>({revision:2})});
 await failed.ctx.commitReaderSettings();
 assert.equal(failed.writes.has('newsgo.account.unsynced.v1.reader'),false,'retry marker cleared after success');
 assert.equal(failed.controls['reader-save-changes'].disabled,true);
 let quota=fixture({role:'user',failStorage:true});
 quota.ctx.stageReaderSetting('newsgo.maximumLeadStories.v1','4');
 await quota.ctx.commitReaderSettings();
 assert.equal(quota.writes.has('newsgo.maximumLeadStories.v1'),false,'failed browser save rolled back');
 assert.equal(quota.calls.length,0,'failed browser save must not PUT');
 assert.equal(quota.controls['reader-save-changes'].disabled,false);
 assert.match(quota.controls['reader-save-status'].textContent,/Save failed/);
 // Reading source selection is navigation history, not a staged preference.
 // Execute the shipped function to check set/remove and storage-error behavior.
 const navigationStart=script.indexOf('function rememberReadingSelection(value)');
 const navigationEnd=script.indexOf('function message(t)',navigationStart);
 assert.ok(navigationStart>=0&&navigationEnd>navigationStart,'navigation function exists');
 const navigationSource=script.slice(navigationStart,navigationEnd);
 const navigationWrites=new Map(),navigationStages=[];
 const navigationContext={
  readingSelectionKey:'newsgo.lastReadingSelection.v1',
  localStorage:{
   setItem:(key,value)=>navigationWrites.set(key,value),
   removeItem:key=>navigationWrites.delete(key)
  },
  stageReaderSetting:(...args)=>navigationStages.push(args),
  String
 };
 vm.createContext(navigationContext);
 vm.runInContext(navigationSource,navigationContext);
 navigationContext.rememberReadingSelection('https://example.org/feed');
 assert.equal(navigationWrites.get(navigationContext.readingSelectionKey),'https://example.org/feed');
 assert.equal(navigationStages.length,0,'selecting Read Articles never stages unsaved edits');
 navigationContext.rememberReadingSelection(null);
 assert.equal(navigationWrites.has(navigationContext.readingSelectionKey),false);
 assert.equal(navigationStages.length,0,'clearing a removed reading selection never stages edits');
 navigationContext.localStorage.setItem=()=>{throw Error('storage unavailable')};
 assert.doesNotThrow(()=>navigationContext.rememberReadingSelection('all'),'failed storage must not block reading');
 assert.ok(script.includes('activeReadingSelection=url;rememberReadingSelection(url)'), 'source selection uses non-staging path');
 assert.ok(script.includes("activeReadingSelection='all';rememberReadingSelection('all')"), 'all subscriptions uses non-staging path');
 assert.ok(!script.includes('stageReaderSetting(readingSelectionKey'),'no other reading navigation stage calls');
 // A restored reading selection must wait for the Homepage verifier.
 // Execute the real refresh() with no verified identity; no API call is allowed.
 const refreshBegin=script.indexOf('async function refresh(){');
 const refreshEnd=script.indexOf('// Only show safe, identifying feed URL components',refreshBegin);
 assert.ok(refreshBegin>=0&&refreshEnd>refreshBegin,'feed refresh function found');
 const refreshText=script.slice(refreshBegin,refreshEnd);
 let fetchCalls=0,renderCalls=0;
 const loading={textContent:'',classList:{remove(){},add(){}}};
 const refreshContext={
  performance:{now:()=>0},
  refreshSequence:0,
  scraperTraceBatch:[],
  updateScraperTraceExportVisibility(){},
  el:id=>id==='personal-loading'?loading:null,
  activeReadingSelection:'all',
  homepageIdentity:null,
  suiteHandoff:'a'.repeat(64),
  suiteExpiresAt:Date.now()+60000,
  rejectedLegacyLocalMode:false,
  refreshPendingAfterAuth:false,
  renderReading(){renderCalls++},
  renderPersonalArticles(){},
  feeds:[],
  Date,Promise,
  authorizedNewsgoFetch(){fetchCalls++;throw Error('should not fetch before verification')}
 };
 vm.createContext(refreshContext);
 vm.runInContext(refreshText,refreshContext);
 await refreshContext.refresh();
 assert.equal(refreshContext.refreshPendingAfterAuth,true,'initial refresh queued until verified');
 assert.match(loading.textContent,/Waiting for Homepage/);
 assert.equal(fetchCalls,0,'no pre-auth API calls');
 assert.equal(renderCalls,0,'no feed rendering required before auth');
 refreshContext.homepageIdentity={role:'guest',user:'guest-001'};
 await refreshContext.refresh();
 assert.equal(refreshContext.refreshPendingAfterAuth,false,'pending refresh clears on authenticated attempt');
 assert.equal(fetchCalls,0,'no feeds means no network request');
 assert.equal(renderCalls,1,'verified refresh resumes normal reading flow');
 // Verify the real session-success branch consumes the pending request.
 assert.ok(script.includes('if(refreshPendingAfterAuth){refreshPendingAfterAuth=false;void refresh()}'),
   'successful Homepage profile verification retries a queued feed refresh');
 assert.ok(script.indexOf('if(refreshPendingAfterAuth){refreshPendingAfterAuth=false;void refresh()}')>
   script.indexOf('  allowVerifiedNewsgo();',script.indexOf('async function checkHomepageSession')),
   'resume occurs only after the browser is marked verified');
 assert.ok(script.includes("if(adminDirty()||hasUnsavedReaderSettings())"),'Safari unload safeguard');
 assert.ok(!script.includes('scheduleAccountSave('),'no automatic account synchronization remains');
 assert.ok(script.includes('localStorage.setItem(field.key'),'recent searches remain browser operational history');
 assert.ok(script.includes('localStorage.removeItem(field.key'),'history clearing remains immediate');
 assert.ok(script.includes("authorizedNewsgoFetch(imageURL,{signal:controller.signal,cache:'no-store'})"),'protected thumbnail requests use the verified NewsGo fetch path');
 assert.ok(script.includes('URL.createObjectURL(blob)')&&script.includes('URL.revokeObjectURL(src)'),'thumbnail blob URL lifecycle is bounded');
 assert.ok(script.includes('for(const controller of thumbnailRequests)controller.abort()'),'stale thumbnail requests are aborted');
 assert.ok(script.includes("if(rssMetadataOnly){"),'RSS image rendering remains a separate path');
 assert.ok(!script.includes('img.src=a.image;\\n   const imageStart'),'normal images no longer use a naked ScraperGo img src');
 console.log('NewsGo Save Changes: staged, personal, Guest, retry, rollback and unload checks passed.');
}
run().catch(e=>{console.error(e);process.exitCode=1});
