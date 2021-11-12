export function makeViral(tweet: string): string {
  return tweet.split(' ').join(' 👏 ')
}
