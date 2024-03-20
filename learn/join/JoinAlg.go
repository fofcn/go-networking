package join

import "sort"

// 用户表
type TradeUser struct {
	Id   int
	Name string
}

// 订单表
type TradeOrder struct {
	Id     int
	UserId int
	Amount uint
}

// 结果表
// Join条件: user.id=order.user_id
type UserOrderView struct {
	UserId      int
	UserName    string
	OrderId     int
	OrderAmount uint
}

// 两个循环进行联接查询
// 算法复杂度：O(M * N)
// 假设用户表有M条记录， 订单表有N条记录
func NestedLoopJoin(users []TradeUser, orders []TradeOrder) []*UserOrderView {
	var userOrderViews []*UserOrderView = make([]*UserOrderView, 0)

	// 遍历用户表
	for _, user := range users {
		// 遍历订单表
		for _, order := range orders {
			// 条件匹配
			if user.Id == order.UserId {
				// Join条件满足添加视图结果
				userOrderViews = append(userOrderViews, &UserOrderView{
					UserId:      user.Id,
					UserName:    user.Name,
					OrderId:     order.Id,
					OrderAmount: order.Amount,
				})
			}
		}
	}

	return userOrderViews
}

// 哈希进行联接查询
// 算法复杂度：O(M + N)
// 假设用户表有M条记录， 订单表有N条记录
func HashJoin(users []TradeUser, orders []TradeOrder) []*UserOrderView {
	var userOrderViews []*UserOrderView = make([]*UserOrderView, 0)

	// 将用户表以用户ID为Key，用户为Value转换为Hash表
	// 算法复杂度：O(M)
	userTable := make(map[int]TradeUser)
	for _, user := range users {
		userTable[user.Id] = user
	}

	// 遍历订单表，查找用户
	// 算法复杂度：O(N)
	for _, order := range orders {
		// 复杂度，接近:O(1)
		if user, exists := userTable[order.UserId]; exists {
			// Join条件满足添加视图结果
			userOrderViews = append(userOrderViews, &UserOrderView{
				UserId:      user.Id,
				UserName:    user.Name,
				OrderId:     order.Id,
				OrderAmount: order.Amount,
			})
		}
	}

	return userOrderViews
}

// 排序进行联接查询
// 算法复杂度：O(M log M + N log N)
// 假设用户表有M条记录， 订单表有N条记录
func SortJoin(users []TradeUser, orders []TradeOrder) []*UserOrderView {
	var userOrderViews []*UserOrderView = make([]*UserOrderView, 0)

	// 排序user表
	// 算法复杂度：O(M log M)
	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})

	// 排序order表
	// 算法复杂度：O(N log N)
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Id < orders[j].Id
	})

	// 遍历订单表，查找用户
	// 算法复杂度：O(M)
	userIdx := 0
	for _, order := range orders {
		// 在user.id为主键的情况下，这里还可以执行二分查找
		for idx < len(users) && users[userIdx].Id < order.UserId {
			userIdx++
		}

		// 如果找到用户，添加到结果集合
		if userIdx < len(users) && users[userIdx].id == order.UserId {
			// Join条件满足添加视图结果
			userOrderViews = append(userOrderViews, &UserOrderView{
				UserId:      user.Id,
				UserName:    user.Name,
				OrderId:     order.Id,
				OrderAmount: order.Amount,
			})
		}

	}

	return userOrderViews
}

// 排序进行联接查询
// 算法复杂度：O(M log N)
func IndexJoin(users []TradeUser, orders []TradeOrder) []*UserOrderView {
	userOrderViews := make([]*UserOrderView, 0)

	// 使用已排序的orders作为具有索引的表
	// 复杂度：O(M)
	for _, user := range users {

		// 二分查找来寻找匹配的订单基于user.Id
		// 复杂度：: O(log N)
		index := sort.Search(len(orders), func(i int) bool {
			return orders[i].UserId >= user.Id
		})

		// 检查是否找到订单，以及id是否匹配
		for index < len(orders) && orders[index].UserId == user.Id {
			// Join条件满足添加视图结果
			userOrderViews = append(userOrderViews, &UserOrderView{
				UserId:      user.Id,
				UserName:    user.Name,
				OrderId:     orders[index].OrderId,
				OrderAmount: orders[index].Amount,
			})

			// 继续寻找下一个匹配的订单（如果有多个匹配的订单）
			index++
		}
	}

	return userOrderViews
}
