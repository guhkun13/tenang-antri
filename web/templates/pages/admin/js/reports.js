// Date functions
function setReportDateRange(range) {
    const dateFrom = document.getElementById('dateFrom');
    const dateTo = document.getElementById('dateTo');
    const today = new Date();
    
    switch(range) {
        case 'today':
            dateFrom.value = today.toISOString().split('T')[0];
            dateTo.value = today.toISOString().split('T')[0];
            break;
        case 'yesterday':
            const yesterday = new Date(today.getTime() - 24 * 60 * 60 * 1000);
            dateFrom.value = yesterday.toISOString().split('T')[0];
            dateTo.value = yesterday.toISOString().split('T')[0];
            break;
        case 'week':
            const weekAgo = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000);
            dateFrom.value = weekAgo.toISOString().split('T')[0];
            dateTo.value = today.toISOString().split('T')[0];
            break;
        case 'month':
            const monthAgo = new Date(today.getFullYear(), today.getMonth(), 1);
            dateFrom.value = monthAgo.toISOString().split('T')[0];
            dateTo.value = today.toISOString().split('T')[0];
            break;
        case 'quarter':
            const quarterAgo = new Date(today.getFullYear(), Math.floor(today.getMonth() / 3) * 3, 1);
            dateFrom.value = quarterAgo.toISOString().split('T')[0];
            dateTo.value = today.toISOString().split('T')[0];
            break;
        case 'year':
            const yearAgo = new Date(today.getFullYear() - 1, 0, 1);
            dateFrom.value = yearAgo.toISOString().split('T')[0];
            dateTo.value = today.toISOString().split('T')[0];
            break;
    }
    
    document.getElementById('reportForm').submit();
}

// Tab switching
function showTab(tabName) {
    document.querySelectorAll('.tab-content').forEach(tab => tab.classList.add('hidden'));
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('border-blue-500', 'text-blue-600');
        btn.classList.add('border-transparent', 'text-gray-600');
    });
    
    const tabContent = document.getElementById(tabName + 'TabContent');
    const tabBtn = document.getElementById(tabName + 'Tab');
    if (tabContent) tabContent.classList.remove('hidden');
    if (tabBtn) {
        tabBtn.classList.remove('border-transparent', 'text-gray-600');
        tabBtn.classList.add('border-blue-500', 'text-blue-600');
    }
}

// Export functions
function exportCSV() {
    const dateFrom = document.getElementById('dateFrom').value;
    const dateTo = document.getElementById('dateTo').value;
    const reportType = document.getElementById('reportType').value;
    
    window.open(`/admin/api/export/tickets?date_from=${dateFrom}&date_to=${dateTo}&type=${reportType}`, '_blank');
}

function exportPDF() {
    const dateFrom = document.getElementById('dateFrom').value;
    const dateTo = document.getElementById('dateTo').value;
    const reportType = document.getElementById('reportType').value;
    
    window.open(`/admin/api/export/tickets/pdf?date_from=${dateFrom}&date_to=${dateTo}&type=${reportType}`, '_blank');
}

// Load report data
function loadReportData() {
    const dateFrom = document.getElementById('dateFrom').value;
    const dateTo = document.getElementById('dateTo').value;
    const reportType = document.getElementById('reportType').value;
    
    if (!dateFrom || !dateTo) {
        document.getElementById('noDataState').classList.remove('hidden');
        document.getElementById('reportContent').classList.add('hidden');
        return;
    }
    
    document.getElementById('noDataState').classList.add('hidden');
    document.getElementById('reportContent').classList.remove('hidden');
    document.getElementById('loadingState').classList.remove('hidden');
    
    fetch(`/admin/api/reports/data?date_from=${dateFrom}&date_to=${dateTo}&type=${reportType}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById('loadingState').classList.add('hidden');
            
            if (data.error) {
                document.getElementById('noDataState').classList.remove('hidden');
                document.getElementById('reportContent').classList.add('hidden');
                return;
            }
            
            updateReportStats(data.summary);
            updateCharts(data.summary);
            updateTables(data.details);
        })
        .catch(error => {
            console.error('Gagal memuat data laporan:', error);
            document.getElementById('loadingState').classList.add('hidden');
            document.getElementById('noDataState').classList.remove('hidden');
            document.getElementById('reportContent').classList.add('hidden');
        });
}

function updateReportStats(summary) {
    document.getElementById('totalTickets').textContent = summary.total_tickets || 0;
    document.getElementById('completedTickets').textContent = summary.completed_tickets || 0;
    document.getElementById('noShowTickets').textContent = summary.no_show_tickets || 0;
    document.getElementById('cancelledTickets').textContent = summary.cancelled_tickets || 0;
    document.getElementById('avgWaitTime').textContent = summary.avg_wait_time || 0;
    document.getElementById('avgServiceTime').textContent = summary.avg_service_time || 0;
    document.getElementById('peakHour').textContent = summary.peak_hour || '--';
    document.getElementById('serviceRate').textContent = summary.service_rate || 0;
}

function updateCharts(summary) {
    updateHourlyChart(summary.hourly_distribution);
    updateCategoryChart(summary.category_distribution);
}

function updateHourlyChart(hourlyData) {
    const chartContainer = document.getElementById('hourlyChart');
    if (!hourlyData || !chartContainer) return;
    
    const maxTickets = Math.max(...Object.values(hourlyData));
    const hours = Array.from({length: 17}, (_, i) => i + 8);
    
    chartContainer.innerHTML = hours.map(hour => {
        const ticketCount = hourlyData[hour] || 0;
        const height = maxTickets > 0 ? (ticketCount / maxTickets) * 100 : 0;
        
        return `
            <div class="flex-1 flex items-end">
                <div class="w-full bg-gray-200 rounded-t-lg relative" style="height: ${height}%">
                    <div class="absolute bottom-0 left-0 right-0 bg-blue-500 text-center text-xs text-white py-1">
                        ${ticketCount}
                    </div>
                </div>
                <span class="w-12 text-xs text-gray-600 text-right">${hour.toString().padStart(2, '0')}:00</span>
            </div>
        `;
    }).join('');
}

function updateCategoryChart(categoryData) {
    const chartContainer = document.getElementById('categoryChart');
    if (!categoryData || !chartContainer) return;
    
    chartContainer.innerHTML = Object.entries(categoryData).map(([category, data]) => `
        <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
            <div class="flex items-center">
                <div class="w-4 h-4 rounded-full mr-3" style="background-color: ${data.color_code}"></div>
                <div>
                    <p class="font-medium">${category}</p>
                    <p class="text-sm text-gray-500">${data.prefix}</p>
                </div>
            </div>
            <div class="text-right">
                <div class="text-2xl font-bold" style="color: ${data.color_code}">${data.count}</div>
                <p class="text-sm text-gray-500">tickets</p>
            </div>
        </div>
    `).join('');
}

function updateTables(details) {
    const tableBody = document.getElementById('ticketsTableBody');
    if (!details || !tableBody) return;
    
    tableBody.innerHTML = details.map(ticket => `
        <tr class="hover:bg-gray-50">
            <td class="px-4 py-2">${new Date(ticket.created_at).toLocaleDateString()}</td>
            <td class="px-4 py-2 font-medium">${ticket.ticket_number}</td>
            <td class="px-4 py-2">
                <span class="px-2 py-1 rounded text-xs text-white" style="background-color: ${ticket.category_color_code}">
                    ${ticket.category_prefix}
                </span>
            </td>
            <td class="px-4 py-2">
                <span class="px-2 py-1 rounded-full text-xs font-medium
                      ${ticket.status === 'completed' ? 'bg-green-100 text-green-800' :
                        ticket.status === 'cancelled' ? 'bg-red-100 text-red-800' :
                        ticket.status === 'no_show' ? 'bg-orange-100 text-orange-800' :
                        'bg-blue-100 text-blue-800'}">
                    ${ticket.status}
                </span>
            </td>
            <td class="px-4 py-2">${ticket.wait_time ? Math.floor(ticket.wait_time / 60) + 'm ' + (ticket.wait_time % 60) + 's' : '-'}</td>
            <td class="px-4 py-2">${ticket.service_time ? Math.floor(ticket.service_time / 60) + 'm ' + (ticket.service_time % 60) + 's' : '-'}</td>
        </tr>
    `).join('');
}

// Auto-load data when page loads
document.addEventListener('DOMContentLoaded', function() {
    loadReportData();
});

// Handle form submission
document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('reportForm');
    if (form) {
        form.addEventListener('submit', function(e) {
            e.preventDefault();
            loadReportData();
        });
    }
});
