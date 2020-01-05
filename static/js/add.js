let addQuoteEndpoint = "http://localhost:8080/profile/add";

function addHandler() {
    let tmp = document.getElementById("add-form");

    document.getElementById("button-id").addEventListener("click", function(event) {
        let x = document.getElementById("add-form").elements[0].value;
        let y = document.getElementById("add-form").elements[1].value;
        if (!x || !y)
            return;
        addNewQuote(x, y);
    });
}

function addNewQuote(data, author) {
    let req = new XMLHttpRequest();

    req.onreadystatechange = function () {
        if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
            alert("Successful add");
            window.location.replace("http://localhost:8080/login");
        }
    };

    let json = {
        data: data,
        author: author
    };

    req.open("POST", addQuoteEndpoint);
    req.send(JSON.stringify(json));
}

addHandler();
