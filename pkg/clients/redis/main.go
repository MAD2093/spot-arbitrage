package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/MAD2093/shared-go/pkg/models/spotarbitrage"
	"mad-scanner.com/scriner/common"

	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
)

// слушает common.ServerData канал и отправляет данные в redis
func RedisChannelListener(channel chan common.ServerData) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "your_strong_password",
		DB:       0,
	})
	ctx := context.Background()

	// Проверка подключения
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	// Инициализация RedisJSON
	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClientWithContext(ctx, rdb)

	// Проверка и создание вложенных структур, если не существуют
	ensurePathExists := func(path string) {
		_, err := rh.JSONGet("arbitrage:spot", path)
		if err != nil {
			_, setErr := rh.JSONSet("arbitrage:spot", path, map[string]interface{}{})
			if setErr != nil {
				fmt.Printf("Ошибка при создании пути %s: %v\n", path, setErr)
			}
		}
	}
	//
	sanitize := func(s string) string {
		return strings.ReplaceAll(s, ".", "_")
	}

	// Удаляем старое дерево полностью
	_, err = rdb.Del(ctx, "arbitrage:spot").Result()
	if err != nil {
		log.Fatalf("Ошибка при удалении ключа Redis: %v", err)
	}
	fmt.Println("Ключ 'arbitrage:spot' очищен")

	// Создаём новый пустой объект
	_, err = rh.JSONSet("arbitrage:spot", ".", map[string]interface{}{})
	if err != nil {
		log.Fatalf("Ошибка при создании корневого JSON-объекта: %v", err)
	}
	fmt.Println("Ключ 'arbitrage:spot' инициализирован")

	// Обработка входящих данных из канала
	for data := range channel {
		coin := sanitize(data.Symbol)
		from := sanitize(data.WithdrawalExchange)
		to := sanitize(data.DepositExchange)

		ensurePathExists(fmt.Sprintf(".%s", coin))
		ensurePathExists(fmt.Sprintf(".%s.%s", coin, from))

		// Сохраняем данные
		path := fmt.Sprintf(".%s.%s.%s", coin, from, to)
		_, err := rh.JSONSet("arbitrage:spot", path, data)
		if err != nil {
			fmt.Printf("Ошибка сохранения: %v\n", err)
			continue
		}

		// Отправка в Pub
		pubMessage := common.RedisMessage{
			Key: spotarbitrage.RedisRoute{
				Symbol:             data.Symbol,
				WithdrawalExchange: data.WithdrawalExchange,
				DepositExchange:    data.DepositExchange,
			},
			Update: data,
			Type:   "update_data",
		}

		messageBytes, err := json.Marshal(pubMessage)
		if err != nil {
			fmt.Printf("Ошибка сериализации pub/sub сообщения: %v\n", err)
			continue
		}

		err = rdb.Publish(ctx, "arbitrage:spot:update", messageBytes).Err()
		if err != nil {
			fmt.Printf("Ошибка отправки pub/sub сообщения: %v\n", err)
			continue
		}
		err = rdb.Publish(ctx, fmt.Sprintf("arbitrage:spot:%s", coin), messageBytes).Err()
		if err != nil {
			fmt.Printf("Ошибка отправки pub/sub сообщения: %v\n", err)
			continue
		}

		fmt.Println(data.Symbol, data.WithdrawalExchange, data.DepositExchange)
	}
}
