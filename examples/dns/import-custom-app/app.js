function init(config) {}

function onRequest(request, response) {
    response.addCName('foo.example.com');
    response.setTTL(20);
}
