const currentPath = window.location.pathname;
for (const link of document.querySelectorAll('.main-nav a')) {
    if (link.getAttribute('href') === currentPath || (currentPath === '/' && link.getAttribute('href') === '/expenses')) {
        link.classList.add('active');
        link.setAttribute('aria-current', 'page');
    }
}
const today = document.getElementById('today');
if (today) {
    today.textContent = new Date().toLocaleDateString('ru-RU', {day: 'numeric', month: 'long', year: 'numeric'});
}
const money = new Intl.NumberFormat('ru-RU', {style: 'currency', currency: 'RUB', maximumFractionDigits: 2});
let total = 0;
for (const cell of document.querySelectorAll('[data-amount]')) {
    const amount = Number(cell.dataset.amount);
    total += amount;
    cell.textContent = money.format(amount);
}
const totalLabel = document.getElementById('expense-total');
if (totalLabel) {
    totalLabel.textContent = money.format(total);
}
for (const cell of document.querySelectorAll('[data-money]')) {
    cell.textContent = money.format(Number(cell.dataset.money));
}
for (const cell of document.querySelectorAll('[data-date]')) {
    const date = new Date(cell.dataset.date + 'T12:00:00');
    if (!Number.isNaN(date.getTime())) {
        cell.textContent = date.toLocaleDateString('ru-RU', {day: 'numeric', month: 'short', year: 'numeric'});
    }
}
const passwordToggle = document.querySelector('.password-toggle');
if (passwordToggle) {
    passwordToggle.addEventListener('click', function () {
        const input = document.getElementById('password');
        const show = input.type === 'password';
        input.type = show ? 'text' : 'password';
        passwordToggle.textContent = show ? 'Скрыть' : 'Показать';
        passwordToggle.setAttribute('aria-label', show ? 'Скрыть пароль' : 'Показать пароль');
        passwordToggle.setAttribute('aria-pressed', String(show));
    });
}
