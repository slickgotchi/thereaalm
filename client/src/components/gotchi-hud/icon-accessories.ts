export type Emotion = "happy" | "sad" | "mad" | "normal";

export const emotions = [
  "normal", 
  "happy", 
  "sad",
  "mad",
];

export const emotionFrameNumbers: Record<string, number> = {
  normal: 0, // No image for normal state
  happy: 3, // Replace with your happy emoticon (e.g., "/images/happy.png")
  sad: 9,   // Replace with your sad emoticon (e.g., "/images/sad.png")
  mad: 12,   // Replace with your mad emoticon (e.g., "/images/mad.png")
};

export const actionFrameNumbers: Record<string, number> = {
  attack: 0,
  forage: 1,
  chop: 2,
  mine: 3,
  harvest: 4,
  defend: 5,
  rebuild: 6,
  maintain: 7,
  sell: 8,
  buy: 9,
  flee: 10,
  sleep: 11,
  roam: 12
}