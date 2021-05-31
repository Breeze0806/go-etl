//  DBStorage（数据库存储）
//     +--------------------------------------------------------+
//     |                                                        |
//   MasterTable（主表）---------+-------+----------------- SlaveTable（从表）
//     |          	            |       |
//     |                        |       +--compare(比较)--+
//     |                        |      /                  |
//     |                        |     /                   +---TableDiffer(表差异)
//     |				    TableNameMap（表映射关联）                |
//   Tracker（跟踪器）                                         read/writer(读/写)
//     +------------------+                                          |
//     |                  |                                 DifferStorage（差异存储）
// OffsetTracker   PageParamTracker
//  (位移记录器)      (分页记录器)
package plugin
