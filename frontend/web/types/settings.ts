export interface UserProfile {
  id: string;
  name: string;
  email: string;
  avatarUrl?: string;
}

export interface NotificationSettings {
  learningReminder: boolean;
  reviewNotification: boolean;
  emailNotification: boolean;
}

export interface UserSettings {
  profile: UserProfile;
  notifications: NotificationSettings;
  interfaceLanguage: string;
}

export interface Plan {
  type: 'free' | 'premium';
  expiresAt?: string;
}
