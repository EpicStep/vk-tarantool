const urlInput = document.querySelector('#url-input')
const submitButton = document.querySelector('#submit-button')
const congratulationsBtn = document.querySelector('.congratulations-btn')

const congratulationsABtn = document.querySelector('.congratulations-analytics')

const createCongratulation = (text) => {
    congratulationsBtn.textContent = text
    congratulationsBtn.setAttribute("href", text)

    congratulationsABtn.textContent = "Аналитика вашего сокращения"
    let a = new URL(text)
    congratulationsABtn.setAttribute("href", "analytics.html?hash=" + a.pathname)

    document.getElementById("qrcode").innerHTML = ""

    new QRCode(document.getElementById("qrcode"), {
        text: text,
        width: 128,
        height: 128,
        colorDark : "#000000",
        colorLight : "#ffffff",
        correctLevel : QRCode.CorrectLevel.H
    });
}

submitButton.addEventListener('click', () => {
    fetch('/set?url=' + urlInput.value, {
        method: 'POST'
    }).then((res) => {
        if (res.status === 200) {
            res.text().then((data) => {
                createCongratulation(data)
            })
        }
    })
})