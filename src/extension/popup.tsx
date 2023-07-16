export default function Index() {
  async function onClick() {
    console.log("Started")

    chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
      chrome.tabs.sendMessage(tabs[0].id, { message: "ButtonClicked" })
    })
  }

  return <button onClick={onClick}>Click</button>
}
