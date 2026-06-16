window.views = window.views || {};

views.login = {
    render(main) {
        main.innerHTML = `
            <div class="login-container card">
                <h1>Firewall Manager</h1>
                <p class="subtitle">Sign in to manage firewall policies</p>
                <div id="login-error" class="alert alert-error hidden"></div>
                <form id="login-form">
                    <div class="form-group">
                        <label for="username">Username</label>
                        <input type="text" id="username" class="form-control" placeholder="Enter username" required autofocus>
                    </div>
                    <div class="form-group">
                        <label for="password">Password</label>
                        <input type="password" id="password" class="form-control" placeholder="Enter password" required>
                    </div>
                    <button type="submit" class="btn btn-primary" style="width:100%">Sign In</button>
                </form>
            </div>
        `;

        document.getElementById('login-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const errorEl = document.getElementById('login-error');
            errorEl.classList.add('hidden');

            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;

            try {
                await api.login(username, password);
                app.showNav(true);
                app.navigate('#/dashboard');
            } catch (err) {
                errorEl.textContent = err.message || 'Login failed';
                errorEl.classList.remove('hidden');
            }
        });
    }
};
