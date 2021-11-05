const views = document.querySelector('.views-p')

let p_url = location.search.substring(1).split("/")[1]
console.log(p_url)

const getColors = (count) => {
    let arr = []

    for (let k = 0; k < count; k++) {
        let letters = '0123456789ABCDEF'.split('');
        let color = '#';

        for (var i = 0; i < 6; i++ ) {
            color += letters[Math.floor(Math.random() * 16)];
        }

        arr.push(color)
    }

    return arr
}

fetch('/analytics?hash=' + p_url, {
    method: 'GET'
}).then((res) => {
    if (res.status === 200) {
        res.json().then((data) => {
            let osLabels = []
            let osData = []

            for (let key in data["response"]["os"]) {
                osLabels.push(key)
                osData.push(data["response"]["os"][key])
            }

            let browserLabels = []
            let browserData = []

            for (let key in data["response"]["browser"]) {
                browserLabels.push(key)
                browserData.push(data["response"]["browser"][key])
            }

            createAnalyticsOS(osLabels, osData)
            createAnalyticsBrowser(browserLabels, browserData)
            createViews(data["response"]["views"])
        })
    }
})

const createViews = (viewsCount) => {
    views.textContent = "Количество переходов: " + viewsCount
}

const createAnalyticsOS = (labels, data) => {
    // Doughnut chart
    let ctx = document.getElementById('myChart').getContext('2d');
    let myChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                data: data,
                backgroundColor: getColors(data.length),
                borderWidth: 0.5 ,
                borderColor: '#ddd'
            }]
        },
        options: {
            title: {
                display: true,
                text: 'Recommended Daily Diet',
                position: 'top',
                fontSize: 16,
                fontColor: '#111',
                padding: 20
            },
            legend: {
                display: true,
                position: 'bottom',
                labels: {
                    boxWidth: 20,
                    fontColor: '#111',
                    padding: 15
                }
            },
            tooltips: {
                enabled: false
            },
            plugins: {
                datalabels: {
                    color: '#111',
                    textAlign: 'center',
                    font: {
                        lineHeight: 1.6
                    },
                    formatter: function(value, ctx) {
                        return ctx.chart.data.labels[ctx.dataIndex] + '\n' + value + '%';
                    }
                }
            }
        }
    });
}

const createAnalyticsBrowser = (labels, data) => {
    // Doughnut chart
    let ctx = document.getElementById('myChart2').getContext('2d');
    let myChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                data: data,
                backgroundColor: getColors(data.length),
                borderWidth: 0.5 ,
                borderColor: '#ddd'
            }]
        },
        options: {
            title: {
                display: true,
                text: 'Recommended Daily Diet',
                position: 'top',
                fontSize: 16,
                fontColor: '#111',
                padding: 20
            },
            legend: {
                display: true,
                position: 'bottom',
                labels: {
                    boxWidth: 20,
                    fontColor: '#111',
                    padding: 15
                }
            },
            tooltips: {
                enabled: false
            },
            plugins: {
                datalabels: {
                    color: '#111',
                    textAlign: 'center',
                    font: {
                        lineHeight: 1.6
                    },
                    formatter: function(value, ctx) {
                        return ctx.chart.data.labels[ctx.dataIndex] + '\n' + value + '%';
                    }
                }
            }
        }
    });
}