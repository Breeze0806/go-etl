//  DBStorage（数据库存储）
//     +--------------------------------------------+
//     |                                            |
//   MasterTable（主表）  -------+------------+    SlaveTable（从表）
//     |          	            |            |
//     |                        |            +---------DifferStorage（差异存储）
//     |				  TableNameMap（表映射关联）
//   Tracker（跟踪器）
//     +------------------+
//     |                  |
// OffsetTracker     PageParamTracker
//  (位移记录器)       (分页记录器)
package plugin
