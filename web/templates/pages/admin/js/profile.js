async function updateProfile(event) {
    event.preventDefault();
    const form = event.target;
    const data = Object.fromEntries(new FormData(form));
    
    try {
        const response = await fetch('/api/profile', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        
        if (response.ok) {
            alert('Profil berhasil diperbarui');
        } else {
            const error = await response.json();
            alert(error.error || 'Gagal memperbarui profil');
        }
    } catch (error) {
        alert('Network error');
    }
    return false;
}

async function changePassword(event) {
    event.preventDefault();
    const form = event.target;
    const data = Object.fromEntries(new FormData(form));
    
    try {
        const response = await fetch('/api/change-password', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        
        if (response.ok) {
            alert('Kata sandi berhasil diubah');
            form.reset();
        } else {
            const error = await response.json();
            alert(error.error || 'Gagal mengubah kata sandi');
        }
    } catch (error) {
        alert('Network error');
    }
    return false;
}
