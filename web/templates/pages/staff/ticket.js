function openResetModal() {
    document.getElementById('resetModal').classList.remove('hidden');
    document.getElementById('resetModal').classList.add('flex');
}

function closeResetModal() {
    document.getElementById('resetModal').classList.add('hidden');
    document.getElementById('resetModal').classList.remove('flex');
}

function cancelTicket(ticketId) {
    if (!confirm('Apakah Anda yakin ingin membatalkan tiket ini?')) return;

    fetch('/staff/api/tickets/' + ticketId + '/cancel', {
        method: 'POST',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert('Gagal membatalkan tiket: ' + data.error);
        } else {
            location.reload();
        }
    })
    .catch(error => {
        alert('Terjadi kesalahan: ' + error);
    });
}

function resetYesterdayTickets() {
    fetch('/staff/api/tickets/reset-yesterday', {
        method: 'POST',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
    .then(response => response.json())
    .then(data => {
        closeResetModal();
        if (data.error) {
            alert('Gagal mereset tiket: ' + data.error);
        } else {
            alert(data.message);
            location.reload();
        }
    })
    .catch(error => {
        closeResetModal();
        alert('Terjadi kesalahan: ' + error);
    });
}

function viewTicketDetail(ticketId) {
    fetch('/staff/api/tickets/' + ticketId, {
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to load ticket details');
        }
        return response.json();
    })
    .then(ticket => {
        displayTicketDetail(ticket);
        openTicketDetailModal();
    })
    .catch(error => {
        alert('Terjadi kesalahan: ' + error.message);
    });
}

function openTicketDetailModal() {
    document.getElementById('ticketDetailModal').classList.remove('hidden');
    document.getElementById('ticketDetailModal').classList.add('flex');
}

function closeTicketDetailModal() {
    document.getElementById('ticketDetailModal').classList.add('hidden');
    document.getElementById('ticketDetailModal').classList.remove('flex');
}

function displayTicketDetail(ticket) {
    // Set basic info
    document.getElementById('detailTicketNumber').textContent = ticket.ticket_number;
    
    // Set status with color
    const statusEl = document.getElementById('detailStatus');
    let statusText = ticket.status;
    let statusClass = '';
    
    switch(ticket.status) {
        case 'waiting':
            statusText = 'Menunggu';
            statusClass = 'bg-yellow-100 text-yellow-800';
            break;
        case 'serving':
            statusText = 'Melayani';
            statusClass = 'bg-blue-100 text-blue-800';
            break;
        case 'completed':
            statusText = 'Selesai';
            statusClass = 'bg-green-100 text-green-800';
            break;
        case 'no_show':
            statusText = 'Tidak Hadir';
            statusClass = 'bg-orange-100 text-orange-800';
            break;
        case 'cancelled':
            statusText = 'Dibatalkan';
            statusClass = 'bg-red-100 text-red-800';
            break;
    }
    statusEl.textContent = statusText;
    statusEl.className = 'px-2 py-1 rounded-full text-xs font-medium ' + statusClass;
    
    // Set category - Go JSON uses lowercase field names
    document.getElementById('detailCategory').textContent = (ticket.category && ticket.category.name) ? ticket.category.name : '-';
    
    // Set timeline - handle sql.NullTime format {Time: "...", Valid: true/false}
    document.getElementById('detailCreatedAt').textContent = formatDateTime(ticket.created_at);
    document.getElementById('detailCalledAt').textContent = (ticket.called_at && ticket.called_at.Valid) ? formatDateTime(ticket.called_at.Time) : 'Belum dipanggil';
    document.getElementById('detailCompletedAt').textContent = (ticket.completed_at && ticket.completed_at.Valid) ? formatDateTime(ticket.completed_at.Time) : 'Belum selesai';
    
    // Set timing metrics - handle sql.NullInt64 format {Int64: value, Valid: true/false}
    document.getElementById('detailWaitTime').textContent = (ticket.wait_time && ticket.wait_time.Valid) ? formatDuration(ticket.wait_time.Int64) : '-';
    document.getElementById('detailServiceTime').textContent = (ticket.service_time && ticket.service_time.Valid) ? formatDuration(ticket.service_time.Int64) : '-';
    
    // Set counter info if available - Go JSON uses lowercase field names
    const counterSection = document.getElementById('detailCounterSection');
    if (ticket.counter && ticket.counter.number) {
        const counterName = ticket.counter.name || '';
        document.getElementById('detailCounter').textContent = ticket.counter.number + (counterName ? ' - ' + counterName : '');
        counterSection.classList.remove('hidden');
    } else {
        counterSection.classList.add('hidden');
    }
}

function formatDateTime(dateValue) {
    if (!dateValue) return '-';
    
    // Handle sql.NullTime format {Time: "...", Valid: true} or direct string
    let dateString = dateValue;
    if (typeof dateValue === 'object' && dateValue.Time) {
        dateString = dateValue.Time;
    }
    
    const date = new Date(dateString);
    
    // Check if date is valid
    if (isNaN(date.getTime())) {
        return '-';
    }
    
    return date.toLocaleString('id-ID', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function formatDuration(seconds) {
    if (!seconds || seconds === 0) return '0 detik';
    
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    
    let result = [];
    if (hours > 0) result.push(hours + ' jam');
    if (minutes > 0) result.push(minutes + ' menit');
    if (secs > 0 && hours === 0) result.push(secs + ' detik');
    
    return result.join(' ') || '0 detik';
}

function sortBy(column) {
    const sortByInput = document.getElementById('sortByInput');
    const sortOrderInput = document.getElementById('sortOrderInput');
    
    if (sortByInput.value === column) {
        // Toggle order
        sortOrderInput.value = sortOrderInput.value === 'asc' ? 'desc' : 'asc';
    } else {
        // New column, default to desc
        sortByInput.value = column;
        sortOrderInput.value = 'desc';
    }
    
    document.getElementById('filterForm').submit();
}

function setQuickDateRange(range) {
    const today = new Date();
    const dateFromInput = document.querySelector('input[name="date_from"]');
    const dateToInput = document.querySelector('input[name="date_to"]');
    
    const formatDate = (date) => {
        return date.toISOString().split('T')[0];
    };
    
    switch(range) {
        case 'today':
            dateFromInput.value = formatDate(today);
            dateToInput.value = formatDate(today);
            break;
        case 'week':
            const weekStart = new Date(today);
            weekStart.setDate(today.getDate() - today.getDay());
            dateFromInput.value = formatDate(weekStart);
            dateToInput.value = formatDate(today);
            break;
        case 'month':
            const monthStart = new Date(today.getFullYear(), today.getMonth(), 1);
            dateFromInput.value = formatDate(monthStart);
            dateToInput.value = formatDate(today);
            break;
    }
}

function clearFilters() {
    document.querySelector('input[name="date_from"]').value = '';
    document.querySelector('input[name="date_to"]').value = '';
    document.querySelector('select[name="status"]').value = '';
    document.getElementById('sortByInput').value = 'created_at';
    document.getElementById('sortOrderInput').value = 'desc';
    document.getElementById('filterForm').submit();
}

function goToPage(page) {
    const urlParams = new URLSearchParams(window.location.search);
    urlParams.set('page', page);
    window.location.search = urlParams.toString();
}