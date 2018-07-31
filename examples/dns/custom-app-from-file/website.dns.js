var prefixes = [
    'foo',
    'bar',
    'baz'
];

var lastIndex = 0;

function init(config) {}

function onRequest(request, response) {
    if (lastIndex >= prefixes.length) {
        lastIndex = 0;
    }
    var prefix = prefixes[lastIndex++];
    response.addCName(prefix + '.example.com');
    response.setTTL(20);
}
