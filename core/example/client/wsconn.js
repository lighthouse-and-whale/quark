window.onload = () => {
    let at = localStorage.getItem("at");
    let b;
    switch (at) {
        case null:
            break;
        case "":
            break;
        case "undefined":
            break;
        default:
            b = true
    }
    if (b === true) {
        wsConn(at)
    } else {
        // 正式环境这里可以跳转到登录页面
        let url = "http://test:2000/v1.app/accounts/newAccessToken";
        fetch(url).then(res => {
            return res.json()
        }).then(res => {
            console.log(res);
            wsConn(res.data.accessToken);
            localStorage.setItem("at", res.data.accessToken)
        })
    }
};

function wsConn(token) {
    const socket = new WebSocket('ws://test:2000/v1.app/accounts/online?at=' + token);

    socket.onerror = (e) => {
        console.log(e);
    };
    socket.onopen = () => {
        console.log("ws: open");
    };
    socket.onclose = () => {
        console.log('ws: close')
    };
    socket.onmessage = (e) => {
        console.log(`--${e.data}--`);
        switch (e.data) {
            case "200":
                socket.send("200");
                break;
            case "1100":
                localStorage.clear();
                alert("检测到您在别处登录，被迫下线");
                break;
        }
    };
}
