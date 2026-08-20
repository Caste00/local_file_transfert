let allFiles = [];

async function loadFiles() {
    const res = await fetch('/files');
    allFiles = await res.json();
    renderFiles();
}

function renderFiles() {
    const query = document.getElementById('searchInput').value.toLowerCase();

    const filtered = allFiles.filter(name => 
        name.trim().toLowerCase().startsWith(query)
    );

    const list = document.getElementById('fileList');
    list.innerHTML = '';

    filtered.forEach(name => {
        const li = document.createElement('li');

        const link = document.createElement('a');
        link.className = 'filename';
        link.href = '/download/' + encodeURIComponent(name);
        link.textContent = truncateName(name, 65);
        link.title = name;
        link.download = name; // suggerisce il nome file al browser

        const deleteBtn = document.createElement('button');
        deleteBtn.className = 'deleteButton';
        deleteBtn.textContent = "Delete";
        deleteBtn.addEventListener('click', () => deleteFile(name));

        li.appendChild(link);
        li.appendChild(deleteBtn);
        list.appendChild(li);
    });
}

document.getElementById('searchInput').addEventListener('input', renderFiles);

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

async function deleteFile(name) {
    await fetch('/files/' + encodeURIComponent(name), {
        method: 'DELETE'
    });
    loadFiles();
}

function truncateName(name, maxLength) {
    if (name.length <= maxLength) return name;
    return name.slice(0, maxLength) + '...';
}