import setupMock, { successResponseWrap } from '@/utils/setup-mock'
import Mock from 'mockjs'

setupMock({
  setup() {
    // 最新项目
    Mock.mock(new RegExp('/api/user/my-project/list'), () => {
      const contributors = [
        {
          name: '张伟',
          email: 'zhangwei@example.com',
          avatar: '//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/a8c8cdb109cb051163646151a4a5083b.png~tplv-uwbnlip3yd-webp.webp',
        },
        {
          name: '李娜',
          email: 'lina@example.com',
          avatar: '//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/3ee5f13fb09879ecb5185e440cef6eb9.png~tplv-uwbnlip3yd-webp.webp',
        },
        {
          name: '王强',
          email: 'wangqiang@example.com',
          avatar: '//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/8361eeb82904210b4f55fab888fe8416.png~tplv-uwbnlip3yd-webp.webp',
        },
        {
          name: '赵敏',
          email: 'zhaomin@example.com',
          avatar: '//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/a8c8cdb109cb051163646151a4a5083b.png~tplv-uwbnlip3yd-webp.webp',
        },
        {
          name: '刘洋',
          email: 'liuyang@example.com',
          avatar: '//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/3ee5f13fb09879ecb5185e440cef6eb9.png~tplv-uwbnlip3yd-webp.webp',
        },
      ]
      const units = [
        {
          name: '2025年xxx奖学金',
          description: 'Scholarship 1',
        },
        {
          name: '2025年xx奖学金',
          description: 'Scholarship 2',
        },
      ]
      return successResponseWrap(
        units.map((unit, index) => ({
          id: index,
          name: unit.name,
          description: unit.description,
          peopleNumber: Mock.Random.natural(10, 1000),
          contributors,
        }))
      )
    })

    // 最新动态
    Mock.mock(new RegExp('/api/user/latest-activity'), () => {
      const activities = []
      let i = 0
      while (i < 7) {
        if (i % 2 === 0) {
          activities.push({
            id: i,
            title: '提交了材料',
            description: '材料审核通过',
            avatar: '//lf1-xgcdn-tos.pstatp.com/obj/vcloud/vadmin/start.8e0e4855ee346a46ccff8ff3e24db27b.png',
          })
        } else {
          activities.push({
            id: i,
            title: '材料审核通过',
            description: '提交了材料',
            avatar: '//lf1-xgcdn-tos.pstatp.com/obj/vcloud/vadmin/start.8e0e4855ee346a46ccff8ff3e24db27b.png',
          })
        }
        i += 1
      }
      return successResponseWrap(activities)
    })

    Mock.mock(new RegExp('/api/user/my-team/list'), () => {
      return successResponseWrap([
        {
          id: 1,
          avatar: '',
          name: '2024级信息学院',
          peopleNumber: Mock.Random.natural(50, 150),
        },
        {
          id: 2,
          avatar: '',
          name: '2024级软件工程系',
          peopleNumber: Mock.Random.natural(50, 150),
        },
        {
          id: 3,
          avatar: '',
          name: '2024级软工2班',
          peopleNumber: Mock.Random.natural(30, 100),
        },
      ])
    })
  },
})
