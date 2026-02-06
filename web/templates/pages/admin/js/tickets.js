    function openModal(id) {
        const modal = document.getElementById(id);
        if (modal) {
            modal.classList.remove('hidden');
            modal.classList.add('flex');
        }
    }

    function closeModal(id) {
        const modal = document.getElementById(id);
        if (modal) {
            modal.classList.add('hidden');
            modal.classList.remove('flex');
        }
    }

    function addTicket(event) {
        event.preventDefault();
        const form = event.target;
        const formData = new FormData(form);

        fetch('/api/tickets', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams(formData)
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                alert('Error: ' + data.error);
            } else {
                closeModal('addTicketModal');
                form.reset();
                location.reload();
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to create ticket');
        });

        return false;
    }

    function viewTicketDetails(ticketId) {
        fetch('/admin/api/tickets/' + ticketId)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(ticket => {
            console.log('Ticket data:', ticket);
            const detailsDiv = document.getElementById('ticketDetails');
            
            const createdAt = ticket.created_at ? new Date(ticket.created_at).toLocaleString() : '-';
            
            let categoryName = '-';
            if (ticket.category) {
                categoryName = ticket.category.name || '-';
            }
            
            let counterInfo = 'Not assigned';
            if (ticket.counter) {
                counterInfo = ticket.counter.number + ' - ' + ticket.counter.name;
            }
            
            let priorityText = 'Normal';
            if (ticket.priority === 1) priorityText = 'Tinggi';
            else if (ticket.priority === 2) priorityText = 'Segera';
            
            detailsDiv.innerHTML = `
                <div class="grid grid-cols-2 gap-4 mb-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-500">Nomor Tiket</label>
                        <p class="text-lg font-bold text-gray-800">${ticket.ticket_number || '-'}</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-500">Status</label>
                        <p class="text-lg font-bold text-gray-800 capitalize">${ticket.status || '-'}</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-500">Category</label>
                        <p class="text-gray-800">${categoryName}</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-500">Priority</label>
                        <p class="text-gray-800 capitalize">${priorityText}</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-500">Counter</label>
                        <p class="text-gray-800">${counterInfo}</p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-500">Created At</label>
                        <p class="text-gray-800">${createdAt}</p>
                    </div>
                </div>
            `;
            openModal('ticketDetailsModal');
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to load ticket details: ' + error.message);
        });
    }

    function cancelTicket(ticketId) {
        if (!confirm('Apakah Anda yakin ingin membatalkan tiket ini?')) {
            return;
        }

        fetch('/api/tickets/' + ticketId + '/cancel', {
            method: 'POST'
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                alert('Error: ' + data.error);
            } else {
                location.reload();
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to cancel ticket');
        });
    }

    function setQuickDateRange(range) {
        const today = new Date();
        let dateFrom, dateTo;

        switch(range) {
            case 'today':
                dateFrom = today.toISOString().split('T')[0];
                dateTo = today.toISOString().split('T')[0];
                break;
            case 'week':
                const lastWeek = new Date(today);
                lastWeek.setDate(today.getDate() - 7);
                dateFrom = lastWeek.toISOString().split('T')[0];
                dateTo = today.toISOString().split('T')[0];
                break;
            case 'month':
                const lastMonth = new Date(today);
                lastMonth.setMonth(today.getMonth() - 1);
                dateFrom = lastMonth.toISOString().split('T')[0];
                dateTo = today.toISOString().split('T')[0];
                break;
        }

        document.querySelector('input[name="date_from"]').value = dateFrom;
        document.querySelector('input[name="date_to"]').value = dateTo;
        document.getElementById('filterForm').submit();
    }

    function clearFilters() {
        document.querySelector('input[name="search"]').value = '';
        document.querySelector('select[name="status"]').value = '';
        document.querySelector('select[name="category_id"]').value = '';
        document.querySelector('select[name="counter_id"]').value = '';
        document.querySelector('input[name="date_from"]').value = '';
        document.querySelector('input[name="date_to"]').value = '';
        document.getElementById('filterForm').submit();
    }
