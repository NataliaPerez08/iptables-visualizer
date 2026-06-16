window.components = window.components || {};

components.navbar = {
    render() {
        const nav = document.getElementById('navbar');
        if (!nav) return;

        const isAdmin = api.isAdmin();
        const user = api.user;

        nav.innerHTML = `
            <div class="nav-brand">Firewall Manager</div>
            <ul class="nav-links">
                <li><a href="#/dashboard" data-nav>Dashboard</a></li>
                <li><a href="#/policies" data-nav>Policies</a></li>
                ${isAdmin ? '<li><a href="#/audit" data-nav>Audit Log</a></li>' : ''}
            </ul>
            <div class="nav-user">
                <span>${user ? `${user.username} (${user.role})` : ''}</span>
                <button id="logout-btn" class="btn btn-sm btn-outline">Logout</button>
            </div>
        `;

        document.getElementById('logout-btn').addEventListener('click', () => {
            api.logout();
            app.showNav(false);
            app.navigate('#/login');
        });
    },

    highlightActive() {
        const hash = window.location.hash || '#/dashboard';
        document.querySelectorAll('[data-nav]').forEach(l => {
            l.classList.toggle('active', l.getAttribute('href') === hash);
        });
    }
};
