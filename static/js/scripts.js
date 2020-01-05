let getQuotesEndpoint = "http://localhost:8080/profile/get/quotes";
let deleteQuoteEndpoint = "http://localhost:8080/profile/delete";
let updateQuoteEndpoint = "http://localhost:8080/profile/put";

function getQuotes() {
    let req = new XMLHttpRequest();

    req.onreadystatechange = function () {
        if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
            loadQuotes(JSON.parse( this.responseText) );
        }
    };

    req.open("GET", getQuotesEndpoint);
    req.send();
}

function loadQuotes(data) {
    document.getElementById("1").innerHTML = '';

    if (data == null) {
        return
    }

    for (let i = 0; i < data.length; i++) {
        const el = data[i];
        console.log(el);
        addNewArticle(el.data, el.author, el.special);
    }
}


function addNewArticle(data, author, special) {
    let tmp = document.getElementsByTagName("template")[0];
    let articleClone = tmp.content.cloneNode(true);

    articleClone.getElementById("quote-data").textContent = data;
    articleClone.getElementById("quote-author").textContent = author;
    if (special == true) {
        articleClone.querySelector(".quote").className += " special-quote";
    }


    articleClone.getElementById("button-id").addEventListener("click", function(event) {
        let v = event.target.parentElement.nextSibling.nextElementSibling.textContent;
        deleteQuoteByData(v);
        event.stopPropagation();
    });

    articleClone.querySelector(".quote").addEventListener("click", function(event) {
        let v = event.target.parentElement.children[1].textContent;
        console.log(v);
        updateQuote(v);
        event.stopPropagation();
    });

    document.getElementById("1").appendChild(articleClone);
}

function deleteQuoteByData(data) {
    let req = new XMLHttpRequest();

    req.onreadystatechange = function () {
        if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
            getQuotes();
        }
    };

    let json = {
        data: data
    };

    req.open("DELETE", deleteQuoteEndpoint);
    req.send(JSON.stringify(json));
}

function updateQuote(data) {
    let req = new XMLHttpRequest();

    req.onreadystatechange = function () {
        if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
            getQuotes();
        }
    };

    let json = {
        data: data
    };

    req.open("PUT", updateQuoteEndpoint);
    req.send(JSON.stringify(json));
}

getQuotes();