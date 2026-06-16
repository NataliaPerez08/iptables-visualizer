const app = {
    currentView: null,

    async init() {
        if (api.isAuthenticated()) {
            this.showNav(true);
            this.navigate(window.location.hash || '#/dashboard');
        } else {
            this.showNav(false);
            this.navigate('#/login');
        }

        window.addEventListener('hashchange', () => {
            if (!api.isAuthenticated() && window.location.hash !== '#/login') {
                this.navigate('#/login');
                return;
            }
            this.navigate(window.location.hash);
        });

        document.getElementById('logout-btn').addEventListener('click', () => {
            api.logout();
            this.showNav(false);
            this.navigate('#/login');
        });

        document.addEventListener('click', (e) => {
            const navLink = e.target.closest('[data-nav]');
            if (navLink) {
                document.querySelectorAll('[data-nav]').forEach(l => l.classList.remove('active'));
                navLink.classList.add('active');
            }
        });
    },

    navigate(hash) {
        const route = hash || '#/dashboard';
        window.location.hash = route;

        if (!api.isAuthenticated() && route !== '#/login') {
            window.location.hash = '#/login';
            return;
        }

        document.querySelectorAll('[data-nav]').forEach(l => l.classList.remove('active'));
        const activeLink = document.querySelector(`[data-nav][href="${route}"]`);
        if (activeLink) activeLink.classList.add('active');

        this.renderView(route);
    },

    renderView(route) {
        const main = document.getElementById('main-content');
        if (!main) return;

        switch (route) {
            case '#/login':
                views.login.render(main);
                break;
            case '#/dashboard':
                views.dashboard.render(main);
                break;
            case '#/policies':
                views.policies.render(main);
                break;
            case '#/audit':
                views.audit.render(main);
                break;
            default:
                main.innerHTML = '<div class="container"><div class="card"><h2>404 - Page Not Found</h2></div></div>';
        }
    },

    showNav(visible) {
        const nav = document.getElementById('navbar');
        nav.classList.toggle('hidden', !visible);
        if (visible) {
            const info = document.getElementById('user-info');
            if (api.user) {
                info.textContent = `${api.user.username} (${api.user.role})`;
            }
        }
    }
};

document.addEventListener('DOMContentLoaded', () => app.init());
