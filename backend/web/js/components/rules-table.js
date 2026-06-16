window.components.rulesTable = {
    render(rules, options = {}) {
        const { readOnly, onEdit, onDelete, onToggle, showActions } = {
            readOnly: false,
            showActions: true,
            ...options
        };

        if (!rules || rules.length === 0) {
            return '<div class="empty-state"><h3>No rules defined</h3></div>';
        }

        return `
            <div class="rules-table">
                <table>
                    <thead>
                        <tr>
                            <th>#</th>
                            <th>Name</th>
                            <th>Action</th>
                            <th>Protocol</th>
                            <th>Source</th>
                            <th>Destination</th>
                            <th>Port</th>
                            ${showActions && !readOnly ? '<th style="text-align:right">Actions</th>' : ''}
                        </tr>
                    </thead>
                    <tbody>
                        ${rules.map((rule, i) => `
                            <tr class="${rule.enabled !== false ? 'rule-enabled' : 'rule-disabled'}">
                                <td>${i + 1}</td>
                                <td><strong>${rule.name || 'Unnamed'}</strong></td>
                                <td><span class="badge badge-${rule.action === 'drop' || rule.action === 'reject' ? 'failed' : 'active'}">${rule.action}</span></td>
                                <td>${rule.protocol || 'any'}</td>
                                <td>${rule.src_addr || 'any'}${rule.src_port ? ':' + rule.src_port : ''}</td>
                                <td>${rule.dst_addr || 'any'}${rule.dst_port ? ':' + rule.dst_port : ''}</td>
                                <td>${rule.dst_port || '-'}</td>
                                ${showActions && !readOnly ? `
                                    <td style="text-align:right">
                                        <div class="rule-actions">
                                            ${onEdit ? `<button class="btn btn-sm" data-edit="${i}">Edit</button>` : ''}
                                            ${onDelete ? `<button class="btn btn-sm btn-danger" data-delete="${i}">Del</button>` : ''}
                                        </div>
                                    </td>
                                ` : ''}
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            </div>
        `;
    }
};
