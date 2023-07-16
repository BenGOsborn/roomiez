import { useEffect, useState } from "react"

const API_KEY_KEY = "API_KEY"
const SEEN_KEY = "SEEN"

export default function Index() {
  const [key, setKey] = useState<string>("")
  const [seen, setSeen] = useState<string[]>([])

  useEffect(() => {
    const _seen = localStorage.getItem(SEEN_KEY)
    if (_seen) setKey(_seen)
  }, [setSeen])

  useEffect(() => {
    localStorage.setItem(SEEN_KEY, JSON.stringify(seen))
  }, [seen])

  useEffect(() => {
    const _key = localStorage.getItem(API_KEY_KEY)
    if (_key) setKey(_key)
  }, [setKey])

  useEffect(() => {
    localStorage.setItem(API_KEY_KEY, key)
  }, [key])

  async function onClick() {
    chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
      chrome.tabs.sendMessage(tabs[0].id, { key, seen }, (resp) => setSeen(resp.seen))
    })
  }

  return (
    <div>
      <input type="password" value={key} onChange={(e) => setKey(e.target.value)} />
      <button onClick={onClick}>Click</button>
    </div>
  )
}
