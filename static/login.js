const toggleButton = document.getElementById('toggleForm');
const loginForm = document.querySelector('.login-form');
const registrationForm = document.querySelector('.registration-form');
const formPanel = document.querySelector('.form-panel');
const registrationPanel = document.querySelector('.registration-panel');
const panelTitle = document.querySelector('.panel-title');
const subTitle = document.querySelector('.subTitle');

function toggleLoginAndRegistration() {
    if (isRegistrationMode) {
        registrationPanel.style.left = '640px';
        formPanel.style.left = '0';
        toggleButton.innerText = '注册';
        panelTitle.innerText = '还未注册？';
        subTitle.innerText = '立即注册，开启新世界！';
        setTimeout(() => {
            loginForm.style.display = 'flex';
            registrationForm.style.display = 'none';
        }, 300);
    } else {
        registrationPanel.style.left = '0';
        formPanel.style.left = '260px';
        toggleButton.innerText = '登录';
        panelTitle.innerText = '已有账号？';
        subTitle.innerText = '已有账号请登录，欢迎回来！';
        setTimeout(() => {
            loginForm.style.display = 'none';
            registrationForm.style.display = 'flex';
        }, 300);
    }
    isRegistrationMode = !isRegistrationMode
}