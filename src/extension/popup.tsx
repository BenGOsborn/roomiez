import { useEffect, useState } from "react"

const API_KEY_KEY = "API_KEY"
const API_BASE_URL_KEY = "API_BASE_URL"

export default function Index() {
  const [key, setKey] = useState<string>("")
  const [apiBaseUrl, setApiBaseUrl] = useState<string>("")

  const [loaded, setLoaded] = useState<boolean>(false)

  useEffect(() => {
    const _apiBaseUrl = localStorage.getItem(API_BASE_URL_KEY)
    if (_apiBaseUrl) setApiBaseUrl(_apiBaseUrl)

    const _key = localStorage.getItem(API_KEY_KEY)
    if (_key) setKey(_key)

    setLoaded(true)
  }, [setApiBaseUrl, setKey, setLoaded])

  useEffect(() => {
    if (!loaded) return

    localStorage.setItem(API_BASE_URL_KEY, apiBaseUrl)
    localStorage.setItem(API_KEY_KEY, key)
  }, [loaded, apiBaseUrl, key])

  async function onClick() {
    chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
      chrome.tabs.sendMessage(tabs[0].id, { key, apiBaseUrl })
    })
  }

  return (
    <div>
      <input type="text" value={apiBaseUrl} onChange={(e) => setApiBaseUrl(e.target.value)} />
      <input type="password" value={key} onChange={(e) => setKey(e.target.value)} />
      <button onClick={onClick}>Click</button>
    </div>
  )
}
