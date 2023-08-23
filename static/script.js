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
// Custom Kanban Renderer Code
renderer.code = function(code, language) {
  if (language == "kanban") {
    return '<div class="kanban">' + code + '</div>';
  }
};

marked.use({ renderer })

var loc = window.location, new_uri;

let websocketURL = "ws://" + loc.host + "/ws"
webSocket = new WebSocket(websocketURL);
webSocket.onmessage = (event) => {
  console.log(event.data);

  let elem = document.getElementById("content")
  // var currentContent = elem.innerHTML

  let result = event.data
  const regex = /---([\s\S]*?)---/;
  if (result.startsWith("---")) {
    result = result.replace(regex, '');
  }

  elem.innerHTML = marked.parse(result);
  mermaid.run();
}

webSocket.onclose = () => {
  window.close();
};
