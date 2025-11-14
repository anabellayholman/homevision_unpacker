const fileInput = document.getElementById('fileInput');
const drop = document.getElementById('drop');
const results = document.getElementById('results');
const downloadZipBtn = document.getElementById('downloadZip');
let lastFiles = [];

drop.addEventListener('click', ()=>fileInput.click());
drop.addEventListener('dragover', (e)=>{e.preventDefault(); drop.classList.add('bg-blue-50')});
drop.addEventListener('dragleave', ()=>{drop.classList.remove('bg-blue-50')});
drop.addEventListener('drop', async (e)=>{e.preventDefault(); drop.classList.remove('bg-blue-50'); handleFiles(e.dataTransfer.files)});
fileInput.addEventListener('change', ()=>handleFiles(fileInput.files));

function detectMagic(buf){
const signatures = {
'jpg': new Uint8Array([0xFF,0xD8,0xFF]),
'png': new Uint8Array([0x89,0x50,0x4E,0x47]),
'gif': new Uint8Array([0x47,0x49,0x46,0x38]),
'pdf': new Uint8Array([0x25,0x50,0x44,0x46]),
'zip': new Uint8Array([0x50,0x4B,0x03,0x04]),
'xml': new TextEncoder().encode('<?xml')
};
for(const k in signatures){
const sig = signatures[k];
for(let i=0;i+sig.length<=buf.length;i++){
let ok=true;
for(let j=0;j<sig.length;j++){ if(buf[i+j]!=sig[j]){ok=false;break} }
if(ok) return {type:k, pos:i};
}
}
return {type:null,pos:-1};
}

function parseEnvArrayBuffer(ab){
const data = new Uint8Array(ab);
const sep = new TextEncoder().encode('**%%');
let parts = [];
let idx = 0;
for(let i=0;i<data.length; i++){
let match=true;
for(let j=0;j<sep.length;j++){ if(data[i+j]!=sep[j]){match=false;break} }
if(match){ parts.push(data.slice(idx,i)); i+=sep.length-1; idx = i+1; }
}
parts.push(data.slice(idx));
let files = [];
for(let p of parts){
let trimmed = trimZeros(p);
if(trimmed.length==0) continue;
let det = detectMagic(trimmed);
let header, content;
if(det.pos>=0){ header = trimmed.slice(0, det.pos); content = trimmed.slice(det.pos); }
else{
let sigIndex = indexOfSequence(trimmed, new TextEncoder().encode('_SIG/'));
if(sigIndex>=0){ header = trimmed.slice(0,sigIndex); content = trimmed.slice(sigIndex+5); }
else { continue; }
}
let meta = parseHeader(header);
let name = meta['FILENAME'] || ('block-'+Math.random().toString(36).substring(2,8));
let ext = meta['EXT'] || det.type || '';
files.push({name, ext, data:content});
}
return files;
}

function trimZeros(u8){
let start=0; while(start<u8.length && u8[start]==0)start++;
let end=u8.length; while(end>start && u8[end-1]==0) end--;
return u8.slice(start,end);
}

function indexOfSequence(buf, seq){
for(let i=0;i+seq.length<=buf.length;i++){
let ok=true; for(let j=0;j<seq.length;j++){ if(buf[i+j]!=seq[j]){ok=false;break} } if(ok) return i;
}
return -1;
}

function parseHeader(u8){
let s = new TextDecoder().decode(u8);
let lines = s.split('\n');
let obj = {};
for(let ln of lines){
ln = ln.trim(); if(!ln) continue;
let parts = ln.split('/'); if(parts.length>=2) obj[parts[0].trim()] = parts.slice(1).join('/').trim();
}
return obj;
}

async function handleFiles(list){
let f = list[0];
if(!f) return;
const ab = await f.arrayBuffer();
const files = parseEnvArrayBuffer(ab);
lastFiles = files;
renderFiles(files);
}

function renderFiles(files){
results.innerHTML = '';
if(files.length==0){ results.textContent = 'No files found'; return; }
let table = document.createElement('div');
table.className='bg-white shadow rounded p-4';
for(let fi of files){
let el = document.createElement('div');
el.className='border-b py-2 flex justify-between items-center';
let left = document.createElement('div');
left.innerHTML = '<strong>'+escapeHtml(fi.name)+'</strong>.'+escapeHtml(fi.ext)+' <span class="text-sm text-gray-500">('+fi.data.length+' bytes)</span>';
let btn = document.createElement('button');
btn.className='px-3 py-1 bg-green-600 text-white rounded';
btn.textContent='Download';
btn.onclick = ()=>downloadFile(fi);
el.appendChild(left); el.appendChild(btn); table.appendChild(el);
}
results.appendChild(table);
}

function escapeHtml(s){ return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }

function downloadFile(fi){
const blob = new Blob([fi.data], {type:'application/octet-stream'});
const a = document.createElement('a');
a.href = URL.createObjectURL(blob);
let filename = fi.name;
if(fi.ext && !filename.endsWith('.'+fi.ext)) filename = filename + '.' + fi.ext;
a.download = filename;
document.body.appendChild(a);
a.click();
a.remove();
}

downloadZipBtn.addEventListener('click', ()=>{
if(!lastFiles || lastFiles.length==0) return;
let zip = new JSZip();
for(let f of lastFiles){
let name = f.name; if(f.ext && !name.endsWith('.'+f.ext)) name = name + '.' + f.ext;
zip.file(name, f.data);
}
zip.generateAsync({type:'blob'}).then(function(content){ const a=document.createElement('a'); a.href=URL.createObjectURL(content); a.download='extracted_files.zip'; a.click(); });
});
