window.views.policies = {
    editingId: null,

    async render(main) {
        main.innerHTML = '<div class="container"><h2>Policies</h2><div class="loading">Loading</div></div>';

        try {
            const policies = await api.get('/policies');

            main.innerHTML = `
                <div class="container">
                    <div class="flex-between mb-16">
                        <h2>Firewall Policies</h2>
                        <button id="new-policy-btn" class="btn btn-primary">+ New Policy</button>
                    </div>
                    <div id="policy-list">
                        ${policies.length === 0 ? `
                            <div class="card empty-state">
                                <h3>No policies defined</h3>
                                <p>Create your first firewall policy to start managing rules</p>
                            </div>
                        ` : `
                        <div class="card">
                            <table>
                                <thead>
                                    <tr>
                                        <th>Name</th>
                                        <th>Status</th>
                                        <th>Rules</th>
                                        <th>Version</th>
                                        <th>Updated</th>
                                        <th>Actions</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    ${policies.map(p => `
                                        <tr>
                                            <td><strong>${p.name}</strong></td>
                                            <td><span class="badge badge-${p.status}">${p.status}</span></td>
                                            <td>${p.rules ? p.rules.length : 0}</td>
                                            <td>v${p.version}</td>
                                            <td>${new Date(p.updated_at).toLocaleDateString()}</td>
                                            <td>
                                                <div class="rule-actions">
                                                    <button class="btn btn-sm" onclick="window.views.policies.edit(${p.id})">Edit</button>
                                                    ${api.isEditor() ? `
                                                        <button class="btn btn-sm btn-success" onclick="window.views.policies.apply(${p.id})">Apply</button>
                                                    ` : ''}
                                                    <button class="btn btn-sm btn-outline" onclick="window.views.policies.dryRun(${p.id})">Dry-Run</button>
                                                    ${api.isAdmin() ? `
                                                        <button class="btn btn-sm btn-danger" onclick="window.views.policies.remove(${p.id})">Delete</button>
                                                    ` : ''}
                                                </div>
                                            </td>
                                        </tr>
                                    `).join('')}
                                </tbody>
                            </table>
                        </div>`}
                    </div>
                    <div id="policy-editor"></div>
                </div>
            `;

            document.getElementById('new-policy-btn').addEventListener('click', () => this.showEditor());
        } catch (err) {
            main.innerHTML = `<div class="container"><div class="alert alert-error">${err.message}</div></div>`;
        }
    },

    async edit(id) {
        try {
            const policy = await api.get(`/policies/${id}`);
            this.editingId = id;
            this.showEditor(policy);
        } catch (err) {
            alert('Failed to load policy: ' + err.message);
        }
    },

    async remove(id) {
        if (!confirm('Delete this policy permanently?')) return;
        try {
            await api.del(`/policies/${id}`);
            app.navigate('#/policies');
        } catch (err) {
            alert('Failed to delete: ' + err.message);
        }
    },

    async apply(id) {
        if (!confirm('Apply this policy to the firewall?')) return;
        try {
            const result = await api.post(`/policies/${id}/apply`);
            alert(`Policy applied successfully.\n${result.count} rules deployed.`);
            app.navigate('#/policies');
        } catch (err) {
            alert('Failed to apply: ' + err.message);
        }
    },

    async dryRun(id) {
        try {
            const policy = await api.get(`/policies/${id}`);
            const result = await api.post(`/policies/${id}/dry-run`);
            const output = result.rules.join('\n');
            const valid = policy.status !== 'failed';

            const modal = document.createElement('div');
            modal.className = 'modal-overlay';
            modal.innerHTML = `
                <div class="modal">
                    <div class="modal-header">
                        <h3>Dry-Run: ${policy.name}</h3>
                        <button class="modal-close" onclick="this.closest('.modal-overlay').remove()">&times;</button>
                    </div>
                    <div class="form-group">
                        <label>Validation</label>
                        <span class="badge badge-${valid ? 'active' : 'failed'}">${valid ? 'Valid' : 'Issues Found'}</span>
                    </div>
                    <div class="form-group">
                        <label>Driver</label>
                        <span>${result.driver}</span>
                    </div>
                    <div class="form-group">
                        <label>Generated Commands (${result.count})</label>
                        <div class="dry-run-output">${output}</div>
                    </div>
                    <div class="mt-16">
                        <button class="btn btn-outline" onclick="this.closest('.modal-overlay').remove()">Close</button>
                    </div>
                </div>
            `;
            document.body.appendChild(modal);
            modal.addEventListener('click', (e) => { if (e.target === modal) modal.remove(); });
        } catch (err) {
            alert('Dry-run failed: ' + err.message);
        }
    },

    showEditor(policy) {
        const editor = document.getElementById('policy-editor');
        const isEdit = !!policy;

        const rules = policy ? policy.rules : [{ id: crypto.randomUUID(), order: 1, name: '', action: 'accept', protocol: 'any', src_addr: '0.0.0.0/0', dst_addr: '0.0.0.0/0', enabled: true }];

        editor.innerHTML = `
            <div class="modal-overlay">
                <div class="modal">
                    <div class="modal-header">
                        <h3>${isEdit ? 'Edit Policy' : 'New Policy'}</h3>
                        <button class="modal-close" onclick="window.views.policies.closeEditor()">&times;</button>
                    </div>
                    <form id="policy-form">
                        <div class="form-group">
                            <label for="policy-name">Policy Name</label>
                            <input type="text" id="policy-name" class="form-control" value="${policy ? policy.name : ''}" required>
                        </div>
                        <div class="form-group">
                            <label for="policy-desc">Description</label>
                            <textarea id="policy-desc" class="form-control">${policy ? policy.description : ''}</textarea>
                        </div>

                        <div class="flex-between">
                            <h4>Rules</h4>
                            <button type="button" class="btn btn-sm" onclick="window.views.policies.addRule()">+ Add Rule</button>
                        </div>
                        <div id="rules-container">
                            ${rules.map((rule, i) => window.views.policies.renderRule(rule, i)).join('')}
                        </div>

                        <div class="flex gap-16 mt-16">
                            <button type="submit" class="btn btn-primary">${isEdit ? 'Update' : 'Create'} Policy</button>
                            <button type="button" class="btn btn-outline" onclick="window.views.policies.closeEditor()">Cancel</button>
                        </div>
                    </form>
                </div>
            </div>
        `;

        document.getElementById('policy-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const data = window.views.policies.collectForm();

            if (!data.name) { alert('Policy name is required'); return; }

            try {
                if (isEdit) {
                    await api.put(`/policies/${policy.id}`, data);
                } else {
                    await api.post('/policies', data);
                }
                window.views.policies.closeEditor();
                app.navigate('#/policies');
            } catch (err) {
                alert('Failed to save policy: ' + err.message);
            }
        });
    },

    renderRule(rule, index) {
        const actions = ['accept', 'drop', 'reject', 'log'];
        const protocols = ['any', 'tcp', 'udp', 'icmp'];

        return `
            <div class="rule-form-section" data-rule-index="${index}">
                <div class="rule-form-header">
                    <h4>Rule #${index + 1}</h4>
                    <div class="flex gap-8">
                        <label class="flex gap-8" style="align-items:center;font-size:0.85rem">
                            <input type="checkbox" onchange="window.views.policies.toggleRule(${index})" ${rule.enabled !== false ? 'checked' : ''}>
                            Enabled
                        </label>
                        <button type="button" class="btn btn-sm btn-danger" onclick="window.views.policies.removeRule(${index})">Remove</button>
                    </div>
                </div>
                <div class="grid-3">
                    <div class="form-group">
                        <label>Rule Name</label>
                        <input type="text" class="form-control rule-name" value="${rule.name || ''}" placeholder="e.g. Allow HTTP">
                    </div>
                    <div class="form-group">
                        <label>Action</label>
                        <select class="form-control rule-action">
                            ${actions.map(a => `<option value="${a}" ${rule.action === a ? 'selected' : ''}>${a.toUpperCase()}</option>`).join('')}
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Protocol</label>
                        <select class="form-control rule-protocol">
                            ${protocols.map(p => `<option value="${p}" ${rule.protocol === p ? 'selected' : ''}>${p.toUpperCase()}</option>`).join('')}
                        </select>
                    </div>
                </div>
                <div class="grid-2">
                    <div class="form-group">
                        <label>Source Address</label>
                        <input type="text" class="form-control rule-src" value="${rule.src_addr || '0.0.0.0/0'}" placeholder="0.0.0.0/0">
                    </div>
                    <div class="form-group">
                        <label>Source Port</label>
                        <input type="text" class="form-control rule-src-port" value="${rule.src_port || ''}" placeholder="e.g. 1024:65535">
                    </div>
                </div>
                <div class="grid-2">
                    <div class="form-group">
                        <label>Destination Address</label>
                        <input type="text" class="form-control rule-dst" value="${rule.dst_addr || '0.0.0.0/0'}" placeholder="0.0.0.0/0">
                    </div>
                    <div class="form-group">
                        <label>Destination Port</label>
                        <input type="text" class="form-control rule-dst-port" value="${rule.dst_port || ''}" placeholder="e.g. 80,443">
                    </div>
                </div>
                <div class="grid-3">
                    <div class="form-group">
                        <label>In Interface</label>
                        <input type="text" class="form-control rule-in-iface" value="${rule.in_interface || ''}" placeholder="eth0">
                    </div>
                    <div class="form-group">
                        <label>Out Interface</label>
                        <input type="text" class="form-control rule-out-iface" value="${rule.out_interface || ''}" placeholder="eth1">
                    </div>
                    <div class="form-group">
                        <label>State</label>
                        <input type="text" class="form-control rule-state" value="${rule.state || ''}" placeholder="NEW,ESTABLISHED">
                    </div>
                </div>
                <div class="form-group">
                    <label>Log Prefix</label>
                    <input type="text" class="form-control rule-log-prefix" value="${rule.log_prefix || ''}" placeholder="FW-DROP: ">
                </div>
                <input type="hidden" class="rule-enabled" value="${rule.enabled !== false ? 'true' : 'false'}">
            </div>
        `;
    },

    addRule() {
        const container = document.getElementById('rules-container');
        const count = container.children.length;
        const div = document.createElement('div');
        div.innerHTML = this.renderRule({
            id: crypto.randomUUID(), order: count + 1, name: '', action: 'accept',
            protocol: 'any', src_addr: '0.0.0.0/0', dst_addr: '0.0.0.0/0', enabled: true
        }, count);
        container.appendChild(div.firstElementChild);
    },

    removeRule(index) {
        const el = document.querySelector(`[data-rule-index="${index}"]`);
        if (el) el.remove();
        this.reindexRules();
    },

    toggleRule(index) {
        const el = document.querySelector(`[data-rule-index="${index}"]`);
        if (!el) return;
        const checkbox = el.querySelector('input[type="checkbox"]');
        const hidden = el.querySelector('.rule-enabled');
        hidden.value = checkbox.checked ? 'true' : 'false';
    },

    reindexRules() {
        document.querySelectorAll('[data-rule-index]').forEach((el, i) => {
            el.dataset.ruleIndex = i;
            const h4 = el.querySelector('h4');
            if (h4) h4.textContent = `Rule #${i + 1}`;
        });
    },

    collectForm() {
        const name = document.getElementById('policy-name').value;
        const description = document.getElementById('policy-desc').value;
        const ruleEls = document.querySelectorAll('[data-rule-index]');

        const rules = Array.from(ruleEls).map((el, i) => ({
            id: crypto.randomUUID(),
            order: i + 1,
            name: el.querySelector('.rule-name').value,
            action: el.querySelector('.rule-action').value,
            protocol: el.querySelector('.rule-protocol').value,
            src_addr: el.querySelector('.rule-src').value,
            src_port: el.querySelector('.rule-src-port').value,
            dst_addr: el.querySelector('.rule-dst').value,
            dst_port: el.querySelector('.rule-dst-port').value,
            in_interface: el.querySelector('.rule-in-iface').value,
            out_interface: el.querySelector('.rule-out-iface').value,
            state: el.querySelector('.rule-state').value,
            log_prefix: el.querySelector('.rule-log-prefix').value,
            enabled: el.querySelector('.rule-enabled').value === 'true'
        }));

        return { name, description, rules };
    },

    closeEditor() {
        const overlay = document.querySelector('.modal-overlay');
        if (overlay) overlay.remove();
        this.editingId = null;
    }
};
