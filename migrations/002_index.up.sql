
-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_user_email ON "user"(email);
CREATE INDEX IF NOT EXISTS idx_user_phone ON "user"(phone);
CREATE INDEX IF NOT EXISTS idx_user_photo_user_id ON user_photo(user_id);
CREATE INDEX IF NOT EXISTS idx_swipe_swiper_id ON swipe(swiper_user_id);
CREATE INDEX IF NOT EXISTS idx_swipe_target_id ON swipe(target_user_id);
CREATE INDEX IF NOT EXISTS idx_message_sender_id ON message(sender_id);
CREATE INDEX IF NOT EXISTS idx_message_match_id ON message(match_id);
CREATE INDEX IF NOT EXISTS idx_match_user1_id ON match(user1_id);
CREATE INDEX IF NOT EXISTS idx_match_user2_id ON match(user2_id);
