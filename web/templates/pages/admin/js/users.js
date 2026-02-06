function openModal(id) {
    document.getElementById(id).classList.remove('hidden');
    document.getElementById(id).classList.add('flex');
}

function closeModal(id) {
    document.getElementById(id).classList.add('hidden');
    document.getElementById(id).classList.remove('flex');
}

async function addUser(event) {
    event.preventDefault();
    const form = event.target;
    const data = Object.fromEntries(new FormData(form));
    
    try {
        const response = await fetch('/admin/api/users', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        
        if (response.ok) {
            window.location.reload();
        } else {
            const error = await response.json();
            alert(error.error || 'Gagal menambahkan pengguna');
        }
    } catch (error) {
        alert('Network error');
    }
    return false;
}

async function loadUserForEdit(id) {
    try {
        const response = await fetch(`/admin/api/users/${id}`);
        if (response.ok) {
            const user = await response.json();
            
            document.getElementById('editUserId').value = user.id.Int64;
            document.getElementById('editFullName').value = user.full_name.String || '';
            document.getElementById('editUsername').value = user.username || '';
            document.getElementById('editEmail').value = user.email.String || '';
            document.getElementById('editPhone').value = user.phone.String || '';
            document.getElementById('editRole').value = user.role || 'staff';
            document.getElementById('editCounterId').value = user.counter_id.Int64 || '';
            
            openModal('editUserModal');
        } else {
            alert('Gagal memuat data pengguna');
        }
    } catch (error) {
        alert('Network error');
    }
}

async function updateUser(event) {
    event.preventDefault();
    const form = event.target;
    const id = form.id.value;
    const data = Object.fromEntries(new FormData(form));
    delete data.id;
    
    try {
        const response = await fetch(`/admin/api/users/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        
        if (response.ok) {
            window.location.reload();
        } else {
            const error = await response.json();
            alert(error.error || 'Gagal memperbarui pengguna');
        }
    } catch (error) {
        alert('Network error');
    }
    return false;
}

async function deleteUser(id) {
    if (!confirm('Apakah Anda yakin ingin menghapus pengguna ini?')) return;
    
    try {
        const response = await fetch(`/admin/api/users/${id}`, { method: 'DELETE' });
        if (response.ok) {
            window.location.reload();
        } else {
            alert('Gagal menghapus pengguna');
        }
    } catch (error) {
        alert('Network error');
    }
}

async function resetPassword(id) {
    if (!confirm('Reset kata sandi ke "password123"?')) return;
    
    try {
        const response = await fetch(`/admin/api/users/${id}/reset-password`, { method: 'POST' });
        if (response.ok) {
            const data = await response.json();
            alert(`Kata sandi direset. Kata sandi baru: ${data.password}`);
        } else {
            alert('Gagal mereset kata sandi');
        }
    } catch (error) {
        alert('Network error');
    }
}

function editUser(id) {
    loadUserForEdit(id);
}
