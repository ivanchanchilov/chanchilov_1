async function loadStats() {
    try {
        Chart.defaults.font.family = '-apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif';
        Chart.defaults.font.size = 11;
        Chart.defaults.color = '#89998f';
        Chart.defaults.borderColor = '#eef2ef';
        const from = document.getElementById("from").value;
        const to = document.getElementById("to").value;
        const year = document.getElementById("year").value;
        let response = await fetch("/api/stats/by-category?from=" + encodeURIComponent(from) + "&to=" + encodeURIComponent(to));
        if (!response.ok || response.redirected) {
            throw new Error("Не удалось получить статистику");
        }
        const categories = await response.json();
        const names = [];
        const totals = [];
        for (let i = 0; i < categories.length; i++) {
            names.push(categories[i].name);
            totals.push(categories[i].total);
        }
        if (categories.length === 0) {
            document.getElementById("category-empty").textContent = "За этот период трат нет.";
        }
        new Chart(document.getElementById("category-chart"), {
            type: "pie",
            data: {
                labels: names,
                datasets: [{
                    label: "Сумма трат",
                    data: totals,
                    borderWidth: 3,
                    borderColor: "#ffffff",
                    backgroundColor: ["#9caf9f", "#dcc8a0", "#9db6cc", "#b3a5c7", "#c9a4b6"]
                }]
            },
            options: {responsive: true, maintainAspectRatio: false, plugins: {legend: {position: "bottom", labels: {usePointStyle: true, padding: 20, boxWidth: 7}}}}
        });

        response = await fetch("/api/stats/by-month?year=" + encodeURIComponent(year));
        if (!response.ok || response.redirected) {
            throw new Error("Не удалось получить статистику");
        }
        const months = await response.json();
        const monthTotals = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
        for (let i = 0; i < months.length; i++) {
            monthTotals[months[i].month - 1] = months[i].total;
        }
        if (months.length === 0) {
            document.getElementById("month-empty").textContent = "За этот год трат нет.";
        }
        new Chart(document.getElementById("month-chart"), {
            type: "bar",
            data: {
                labels: ["Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"],
                datasets: [{label: "Сумма трат", data: monthTotals, backgroundColor: "#7da68c", borderRadius: 4, maxBarThickness: 28}]
            },
            options: {responsive: true, maintainAspectRatio: false, plugins: {legend: {display: false}}, scales: {x: {grid: {display: false}}, y: {beginAtZero: true, border: {display: false}}}}
        });
    } catch (error) {
        document.getElementById("stats-error").textContent = "Не удалось загрузить графики. Проверьте подключение к интернету и вход в аккаунт.";
    }
}

loadStats();
