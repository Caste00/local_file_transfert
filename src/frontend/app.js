async function loadFiles() {
    const res = await fetch('/files');
    const names = await res.json();

    const list = document.getElementById('fileList');
    list.innerHTML = '';

    names.forEach(name => {
        const li = document.createElement('li')
        const link = document.createElement('a')
        link.href = '/download/' + encodeURIComponent(name)
        link.textContent = name;
        li.appendChild(link);
        list.appendChild(li);
    });
}

loadFiles();

document.getElementById('uploadForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const input = document.getElementById('fileInput');
    const files = input.files;

    for (const file of files) {
        const formData = new FormData();
        formData.append('file', file);

        await fetch('/upload', {
            method: 'POST',
            body: formData
        });
    }

    input.value = '';
    loadFiles();
})