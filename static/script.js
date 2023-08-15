const { markedHighlight } = globalThis.markedHighlight;
marked.use({
  gfm: true,
  mangle: false,
  headerIds: false,
});


marked.use(markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code, lang) {
        const language = hljs.getLanguage(lang) ? lang : 'plaintext';
        return hljs.highlight(code, { language }).value;
    }
}));
const renderer = new marked.Renderer();
renderer.code = function(code, language) {
  if (language == "mermaid") {
    return '<pre class="mermaid">' + code + '</pre>';
  } else {
    return '<pre><code>' + code + '</code></pre>';
  }
};

marked.use({ renderer })

var loc = window.location, new_uri;

let websocketURL = "ws://" + loc.host + "/ws"
webSocket = new WebSocket(websocketURL);
webSocket.onmessage = (event) => {
  console.log(event.data);
  document.getElementById('content').innerHTML =
    marked.parse(event.data);
  mermaid.run();
}
