export type Emotion = "happy" | "sad" | "mad" | "normal";

export const emotiveAccessories: Record<Emotion, number> = {
  happy: 3, // Replace with your happy emoticon (e.g., "/images/happy.png")
  sad: 9,   // Replace with your sad emoticon (e.g., "/images/sad.png")
  mad: 12,   // Replace with your mad emoticon (e.g., "/images/mad.png")
  normal: 0, // No image for normal state
};