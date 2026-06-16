const API_BASE = '/api/v1';

const api = {
    token: null,
    user: null,

    init() {
        this.token = localStorage.getItem('token');
        const stored = localStorage.getItem('user');
        if (stored) {
            try { this.user = JSON.parse(stored); } catch(e) { this.user = null; }
        }
    },

    async request(method, path, body = null) {
        const headers = { 'Content-Type': 'application/json' };
        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const opts = { method, headers };
        if (body) opts.body = JSON.stringify(body);

        const res = await fetch(`${API_BASE}${path}`, opts);

        if (res.status === 204) return null;

        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.error || `Request failed: ${res.status}`);
        }

        return data;
    },

    get(path) { return this.request('GET', path); },
    post(path, body) { return this.request('POST', path, body); },
    put(path, body) { return this.request('PUT', path, body); },
    del(path) { return this.request('DELETE', path); },

    async login(username, password) {
        const data = await this.post('/auth/login', { username, password });
        this.token = data.token;
        this.user = data.user;
        localStorage.setItem('token', data.token);
        localStorage.setItem('user', JSON.stringify(data.user));
        return data;
    },

    logout() {
        this.token = null;
        this.user = null;
        localStorage.removeItem('token');
        localStorage.removeItem('user');
    },

    isAuthenticated() {
        return !!this.token && !!this.user;
    },

    isAdmin() {
        return this.user && this.user.role === 'admin';
    },

    isEditor() {
        return this.user && (this.user.role === 'admin' || this.user.role === 'editor');
    }
};

api.init();
