import type { PlasmoCSConfig } from "plasmo"

export {}

chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
  const seen = new Set<string>(request.seen)
  const key = request.key

  run(seen, key).then(() => {
    sendResponse({ seen: Array.from(seen) })
  })

  return true
})

async function run(seen: Set<string>, key: string) {
  await new Promise((res) => setTimeout(res, 3000))

  for (const elem of Array.from(document.querySelector("[role=feed]").children)) {
    // Expand the text
    const msg = elem.querySelector("[data-ad-comet-preview=message]")

    if (msg) {
      const more = msg.querySelector("[role=button]")

      if (more) {
        // @ts-ignore
        more.click()
        await new Promise((res) => setTimeout(res, 100))
      }
    }

    // Get the post URL
    const share = elem.querySelector('[aria-label="Send this to friends or post it on your Timeline."]')
    if (!!share) {
      // @ts-ignore
      share.click()
      await new Promise((res) => setTimeout(res, 1000))

      const shareOptionContainer = document.querySelector("[role=dialog]")
      if (!!shareOptionContainer) {
        const shareOptions = shareOptionContainer.querySelectorAll("[role=button]")

        if (!!shareOptions) {
          const copyLink = shareOptions[shareOptions.length - 1]

          copyLink.scrollIntoView()

          const url = await new Promise<string>((res) => {
            copyLink.addEventListener("click", () => {
              navigator.clipboard.readText().then((url) => res(url))
            })
          })

          if (seen.has(url)) break
          seen.add(url)

          //   Post request
          console.log({ url, post: msg.textContent })

          //   **** So now here we can take the values and submit them
        }
      }
    }
  }
}

export const config: PlasmoCSConfig = {
  matches: ["http://www.facebook.com/*"]
}
