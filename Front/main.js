document.getElementById("getBtn").addEventListener("click", getDataFromServer);

function getDataFromServer() {
    let uid = document.getElementById("uid").value
    fetch(`http://127.0.0.1:8080/data/${uid}`, {
        method: "GET",
    }).then(response => {
        return response.json()
    }).then((data) => {
        document.getElementById("msg").textContent = JSON.stringify(data, undefined, 2)
    }).catch((error) => console.error("Error:", error))
}