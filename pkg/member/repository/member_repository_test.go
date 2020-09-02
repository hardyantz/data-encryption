package repository

import (
	"testing"
	"time"

	"github.com/hardyantz/data-encryption/config"
	"github.com/hardyantz/data-encryption/conn"
	"github.com/hardyantz/data-encryption/pkg/member/domain"
)

func BenchmarkMemberRepository_Create(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		member := new(domain.Member)
		db := conn.ConnectDB()
		conf := config.NewConfImpl()

		cacheExpire := 1 * time.Hour

		redis := config.NewCacheRedis(conn.RedisConnect(), cacheExpire)
		memberRepo := NewMemberRepoImplementation(db, conf, redis)
		_ = memberRepo.Create(member)
	}
}
//
//func BenchmarkMemberRepository_GetAll(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		var member domain.Parameters
//		db := conn.ConnectDB()
//		conf := config.NewConfImpl()
//
//		cacheExpire := 1 * time.Hour
//
//		redis := config.NewCacheRedis(conn.RedisConnect(), cacheExpire)
//		memberRepo := NewMemberRepoImplementation(db, conf, redis)
//		_, _ = memberRepo.GetAll(member)
//	}
//}
//
//func BenchmarkMemberRepository_GetOne(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		db := conn.ConnectDB()
//		conf := config.NewConfImpl()
//
//		cacheExpire := 1 * time.Hour
//
//		redis := config.NewCacheRedis(conn.RedisConnect(), cacheExpire)
//		memberRepo := NewMemberRepoImplementation(db, conf, redis)
//		_, _ = memberRepo.GetOne("20200901152143")
//	}
//}
//
//func BenchmarkMemberRepository_Update(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		db := conn.ConnectDB()
//		conf := config.NewConfImpl()
//
//		cacheExpire := 1 * time.Hour
//
//		redis := config.NewCacheRedis(conn.RedisConnect(), cacheExpire)
//		memberRepo := NewMemberRepoImplementation(db, conf, redis)
//		_, _ = memberRepo.GetOne("20200901152143")
//	}
//}
//
//func BenchmarkMemberRepository_Delete(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		db := conn.ConnectDB()
//		conf := config.NewConfImpl()
//
//		cacheExpire := 1 * time.Hour
//
//		redis := config.NewCacheRedis(conn.RedisConnect(), cacheExpire)
//		memberRepo := NewMemberRepoImplementation(db, conf, redis)
//		_ = memberRepo.Delete("20200604213015")
//	}
//}
////
////func BenchmarkMemberRepository_GetOneEmail(b *testing.B) {
////	b.ResetTimer()
////	for i := 0; i < b.N; i++ {
////		db := conn.ConnectDB()
////		conf := config.NewConfImpl()
////
////		cacheExpire := 1 * time.Hour
////
////		redis := config.NewCacheRedis(conn.RedisConnect(), cacheExpire)
////		memberRepo := NewMemberRepoImplementation(db, conf, redis)
////		_, _ = memberRepo.GetOneEmail("customredis@gmail.com")
////	}
////}