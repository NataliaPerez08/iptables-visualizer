window.components.auditLog = {
    render(log) {
        return `
            <div class="audit-entry flex-between">
                <div>
                    <span class="audit-action ${!log.success ? 'failed' : ''}">${log.action}</span>
                    <span>${log.resource}</span>
                    <small class="text-muted">by ${log.username}</small>
                    ${log.details ? `<small class="text-muted">- ${log.details}</small>` : ''}
                </div>
                <small class="text-muted">${new Date(log.created_at).toLocaleString()}</small>
            </div>
        `;
    },

    renderList(logs) {
        if (!logs || logs.length === 0) {
            return '<div class="empty-state"><h3>No audit logs found</h3></div>';
        }

        return `
            <div class="card">
                <table>
                    <thead>
                        <tr>
                            <th>Time</th>
                            <th>User</th>
                            <th>Action</th>
                            <th>Resource</th>
                            <th>Details</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${logs.map(log => `
                            <tr>
                                <td><small>${new Date(log.created_at).toLocaleString()}</small></td>
                                <td>${log.username}</td>
                                <td><span class="audit-action ${!log.success ? 'failed' : ''}">${log.action}</span></td>
                                <td><small>${log.resource}</small></td>
                                <td><small>${log.details || '-'}</small></td>
                                <td>${log.success ? '<span class="badge badge-active">OK</span>' : '<span class="badge badge-failed">FAIL</span>'}</td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            </div>
        `;
    }
};
