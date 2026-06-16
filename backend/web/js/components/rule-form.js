window.components.ruleForm = {
    render(rule = {}, index = 0) {
        const actions = ['accept', 'drop', 'reject', 'log'];
        const protocols = ['any', 'tcp', 'udp', 'icmp'];

        return `
            <div class="rule-form-section" data-rule-index="${index}">
                <div class="rule-form-header">
                    <h4>Rule #${index + 1}</h4>
                    <div class="flex gap-8">
                        <label style="display:flex;align-items:center;gap:6px;font-size:0.85rem">
                            <input type="checkbox" class="rule-enabled-toggle" ${rule.enabled !== false ? 'checked' : ''}>
                            Enabled
                        </label>
                        <button type="button" class="btn btn-sm btn-danger" data-remove-rule>Remove</button>
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

    collect(index) {
        const el = document.querySelector(`[data-rule-index="${index}"]`);
        if (!el) return null;
        return {
            id: crypto.randomUUID(),
            order: index + 1,
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
        };
    }
};
