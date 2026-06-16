window.views.dashboard = {
    async render(main) {
        main.innerHTML = '<div class="container"><h2>Dashboard</h2><div class="loading">Loading</div></div>';

        try {
            const [policies, auditLogs] = await Promise.all([
                api.get('/policies'),
                api.get('/audit?limit=10')
            ]);

            const total = policies.length;
            const active = policies.filter(p => p.status === 'active').length;
            const draft = policies.filter(p => p.status === 'draft').length;
            const failed = policies.filter(p => p.status === 'failed').length;
            const recentLogs = auditLogs || [];

            main.innerHTML = `
                <div class="container">
                    <h2>Dashboard</h2>
                    <div class="stat-cards">
                        <div class="stat-card">
                            <h3>Total Policies</h3>
                            <div class="stat-value stat-primary">${total}</div>
                        </div>
                        <div class="stat-card">
                            <h3>Active</h3>
                            <div class="stat-value stat-success">${active}</div>
                        </div>
                        <div class="stat-card">
                            <h3>Draft</h3>
                            <div class="stat-value stat-warning">${draft}</div>
                        </div>
                        <div class="stat-card">
                            <h3>Failed</h3>
                            <div class="stat-value stat-danger">${failed}</div>
                        </div>
                    </div>

                    <div class="card">
                        <div class="card-header">
                            <span class="card-title">Recent Policies</span>
                            <a href="#/policies" class="btn btn-sm">View All</a>
                        </div>
                        ${policies.length === 0 ? '<div class="empty-state"><h3>No policies yet</h3><p>Create your first firewall policy</p></div>' : `
                        <table>
                            <thead><tr><th>Name</th><th>Status</th><th>Rules</th><th>Version</th><th>Updated</th></tr></thead>
                            <tbody>
                                ${policies.slice(0, 5).map(p => `
                                    <tr onclick="app.navigate('#/policies')" style="cursor:pointer">
                                        <td><strong>${p.name}</strong></td>
                                        <td><span class="badge badge-${p.status}">${p.status}</span></td>
                                        <td>${p.rules ? p.rules.length : 0}</td>
                                        <td>${p.version}</td>
                                        <td>${new Date(p.updated_at).toLocaleDateString()}</td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>`}
                    </div>

                    <div class="card">
                        <div class="card-header">
                            <span class="card-title">Recent Activity</span>
                            <a href="#/audit" class="btn btn-sm">View All</a>
                        </div>
                        ${recentLogs.length === 0 ? '<div class="empty-state"><h3>No recent activity</h3></div>' : `
                        <div>
                            ${recentLogs.map(log => `
                                <div class="audit-entry flex-between">
                                    <div>
                                        <span class="audit-action ${!log.success ? 'failed' : ''}">${log.action}</span>
                                        <span>${log.resource}</span>
                                    </div>
                                    <small class="text-muted">${new Date(log.created_at).toLocaleString()}</small>
                                </div>
                            `).join('')}
                        </div>`}
                    </div>
                </div>
            `;
        } catch (err) {
            main.innerHTML = `<div class="container"><div class="alert alert-error">Failed to load dashboard: ${err.message}</div></div>`;
        }
    }
};
