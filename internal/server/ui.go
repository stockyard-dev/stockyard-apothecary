package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Apothecary</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center}
.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.med{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.med:hover{border-color:var(--leather)}
.med.inactive{opacity:.5}
.med-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.med-name{font-size:.88rem;font-weight:700}
.med-dosage{font-size:.7rem;color:var(--gold);margin-top:.1rem}
.med-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.med-notes{font-size:.65rem;color:var(--cm);margin-top:.3rem;font-style:italic;padding:.3rem .5rem;border-left:2px solid var(--bg3)}
.med-actions{display:flex;gap:.3rem;flex-shrink:0;align-items:center}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.badge.active{border-color:var(--green);color:var(--green)}.badge.inactive{border-color:var(--cm);color:var(--cm)}
.refill-warn{color:var(--red);font-weight:700}
.toggle{position:relative;display:inline-block;width:32px;height:18px}.toggle input{opacity:0;width:0;height:0}
.sl{position:absolute;cursor:pointer;inset:0;background:var(--bg3);transition:.2s;border-radius:18px}
.sl:before{content:'';position:absolute;height:14px;width:14px;left:2px;bottom:2px;background:var(--cm);transition:.2s;border-radius:50%}
.toggle input:checked+.sl{background:var(--green)}.toggle input:checked+.sl:before{transform:translateX(14px);background:var(--cream)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> APOTHECARY</h1><button class="btn btn-p" onclick="openForm()">+ Add Medication</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar"><input class="search" id="search" placeholder="Search medications..." oninput="render()"></div>
<div id="meds"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',meds=[],editId=null;

async function load(){var r=await fetch(A+'/medications').then(function(r){return r.json()});meds=r.medications||[];renderStats();render();}

function renderStats(){
var total=meds.length,active=meds.filter(function(m){return m.active}).length;
var needRefill=meds.filter(function(m){return m.active&&m.refill_date&&new Date(m.refill_date)<=new Date()}).length;
document.getElementById('stats').innerHTML=[
{l:'Medications',v:total},{l:'Active',v:active,c:'var(--green)'},{l:'Need Refill',v:needRefill,c:needRefill>0?'var(--red)':''}
].map(function(x){return '<div class="st"><div class="st-v" style="'+(x.c?'color:'+x.c:'')+'">'+x.v+'</div><div class="st-l">'+x.l+'</div></div>'}).join('');
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var f=meds;
if(q)f=f.filter(function(m){return(m.name||'').toLowerCase().includes(q)||(m.prescriber||'').toLowerCase().includes(q)||(m.pharmacy||'').toLowerCase().includes(q)});
f.sort(function(a,b){return(b.active||0)-(a.active||0)});
if(!f.length){document.getElementById('meds').innerHTML='<div class="empty">No medications tracked.</div>';return;}
var h='';f.forEach(function(m){
var needsRefill=m.active&&m.refill_date&&new Date(m.refill_date)<=new Date();
h+='<div class="med'+(m.active?'':' inactive')+'"><div class="med-top"><div style="flex:1">';
h+='<div class="med-name">'+esc(m.name)+'</div>';
if(m.dosage||m.frequency)h+='<div class="med-dosage">'+esc(m.dosage)+(m.frequency?' &#183; '+esc(m.frequency):'')+'</div>';
h+='</div><div class="med-actions">';
h+='<label class="toggle"><input type="checkbox" '+(m.active?'checked':'')+' onchange="toggleActive(''+m.id+'')"><span class="sl"></span></label>';
h+='<button class="btn btn-sm" onclick="openEdit(''+m.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+m.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
h+='<div class="med-meta">';
h+='<span class="badge '+(m.active?'active':'inactive')+'">'+(m.active?'active':'inactive')+'</span>';
if(m.prescriber)h+='<span>Dr. '+esc(m.prescriber)+'</span>';
if(m.pharmacy)h+='<span>'+esc(m.pharmacy)+'</span>';
if(m.refill_date)h+='<span class="'+(needsRefill?'refill-warn':'')+'">Refill: '+m.refill_date+'</span>';
h+='</div>';
if(m.notes)h+='<div class="med-notes">'+esc(m.notes)+'</div>';
h+='</div>';
});
document.getElementById('meds').innerHTML=h;
}

async function toggleActive(id){var m=null;for(var j=0;j<meds.length;j++){if(meds[j].id===id){m=meds[j];break;}}if(!m)return;
await fetch(A+'/medications/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({active:m.active?0:1})});load();}
async function del(id){if(!confirm('Remove?'))return;await fetch(A+'/medications/'+id,{method:'DELETE'});load();}

function formHTML(med){
var i=med||{name:'',dosage:'',frequency:'',prescriber:'',pharmacy:'',refill_date:'',notes:''};
var isEdit=!!med;
var h='<h2>'+(isEdit?'EDIT MEDICATION':'ADD MEDICATION')+'</h2>';
h+='<div class="fr"><label>Name *</label><input id="f-name" value="'+esc(i.name)+'" placeholder="e.g. Lisinopril"></div>';
h+='<div class="row2"><div class="fr"><label>Dosage</label><input id="f-dosage" value="'+esc(i.dosage)+'" placeholder="e.g. 10mg"></div>';
h+='<div class="fr"><label>Frequency</label><input id="f-freq" value="'+esc(i.frequency)+'" placeholder="e.g. Once daily"></div></div>';
h+='<div class="row2"><div class="fr"><label>Prescriber</label><input id="f-doc" value="'+esc(i.prescriber)+'"></div>';
h+='<div class="fr"><label>Pharmacy</label><input id="f-pharm" value="'+esc(i.pharmacy)+'"></div></div>';
h+='<div class="fr"><label>Next Refill</label><input id="f-refill" type="date" value="'+esc(i.refill_date)+'"></div>';
h+='<div class="fr"><label>Notes</label><input id="f-notes" value="'+esc(i.notes)+'"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button></div>';
return h;
}

function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');document.getElementById('f-name').focus();}
function openEdit(id){var m=null;for(var j=0;j<meds.length;j++){if(meds[j].id===id){m=meds[j];break;}}if(!m)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(m);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}

async function submit(){
var name=document.getElementById('f-name').value.trim();
if(!name){alert('Name is required');return;}
var body={name:name,dosage:document.getElementById('f-dosage').value.trim(),frequency:document.getElementById('f-freq').value.trim(),prescriber:document.getElementById('f-doc').value.trim(),pharmacy:document.getElementById('f-pharm').value.trim(),refill_date:document.getElementById('f-refill').value,notes:document.getElementById('f-notes').value.trim()};
if(editId){await fetch(A+'/medications/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{body.active=1;await fetch(A+'/medications',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();
}

function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
