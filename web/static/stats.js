async function loadStats() {
    try {
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
                    backgroundColor: ["#5b9bd5", "#ed7d31", "#70ad47", "#ffc000", "#9966cc"]
                }]
            }
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
                datasets: [{label: "Сумма трат", data: monthTotals, backgroundColor: "#5b9bd5"}]
            },
            options: {scales: {y: {beginAtZero: true}}}
        });
    } catch (error) {
        document.getElementById("stats-error").textContent = "Не удалось загрузить графики. Проверьте подключение к интернету и вход в аккаунт.";
    }
}

loadStats();
