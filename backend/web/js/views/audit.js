window.views.audit = {
    async render(main) {
        main.innerHTML = `
            <div class="container">
                <h2>Audit Log</h2>
                <div class="card mb-16">
                    <form id="audit-filter-form" class="flex gap-16" style="flex-wrap:wrap;align-items:flex-end">
                        <div class="form-group" style="flex:1;min-width:150px">
                            <label>Action</label>
                            <select id="filter-action" class="form-control">
                                <option value="">All Actions</option>
                                <option value="login">Login</option>
                                <option value="create_policy">Create Policy</option>
                                <option value="update_policy">Update Policy</option>
                                <option value="delete_policy">Delete Policy</option>
                                <option value="apply_policy">Apply Policy</option>
                                <option value="dry_run_policy">Dry-Run</option>
                                <option value="validate_policy">Validate</option>
                                <option value="create_user">Create User</option>
                            </select>
                        </div>
                        <div class="form-group" style="flex:1;min-width:150px">
                            <label>Resource</label>
                            <input type="text" id="filter-resource" class="form-control" placeholder="policy:my-policy">
                        </div>
                        <div class="form-group" style="flex:1;min-width:150px">
                            <label>From</label>
                            <input type="date" id="filter-from" class="form-control">
                        </div>
                        <div class="form-group" style="flex:1;min-width:150px">
                            <label>To</label>
                            <input type="date" id="filter-to" class="form-control">
                        </div>
                        <button type="submit" class="btn btn-primary">Filter</button>
                    </form>
                </div>
                <div id="audit-results"><div class="loading">Loading</div></div>
            </div>
        `;

        await this.loadLogs();

        document.getElementById('audit-filter-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            await this.loadLogs();
        });
    },

    async loadLogs() {
        const resultsEl = document.getElementById('audit-results');
        resultsEl.innerHTML = '<div class="loading">Loading</div>';

        try {
            const action = document.getElementById('filter-action')?.value || '';
            const resource = document.getElementById('filter-resource')?.value || '';
            const from = document.getElementById('filter-from')?.value || '';
            const to = document.getElementById('filter-to')?.value || '';

            let path = '/audit?limit=100';
            if (action) path += `&action=${action}`;
            if (resource) path += `&resource=${resource}`;
            if (from) path += `&from=${new Date(from).toISOString()}`;
            if (to) path += `&to=${new Date(to + 'T23:59:59').toISOString()}`;

            const logs = await api.get(path);

            if (logs.length === 0) {
                resultsEl.innerHTML = '<div class="card empty-state"><h3>No audit logs found</h3></div>';
                return;
            }

            resultsEl.innerHTML = `
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
                    <div class="text-muted mt-16 text-center">${logs.length} entries</div>
                </div>
            `;
        } catch (err) {
            resultsEl.innerHTML = `<div class="alert alert-error">Failed to load audit logs: ${err.message}</div>`;
        }
    }
};
