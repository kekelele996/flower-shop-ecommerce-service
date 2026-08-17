import { useEffect } from 'react';
import { Card, Row, Col, Tag, Typography, Button, message } from 'antd';
import { getMyCoupons, claimCoupon } from '../../api/coupon';
import { useCouponStore } from '../../stores/couponStore';
import { formatPrice, formatPercent } from '../../utils/format';
import EmptyState from '../../components/EmptyState';
import type { CouponVO } from '../../api/coupon';
import { COUPON_TYPE } from '../../constants/enums';

// 我的优惠券：查看与领取。
export default function Coupons() {
  const { coupons, fetchCoupons } = useCouponStore();

  useEffect(() => {
    fetchCoupons();
  }, [fetchCoupons]);

  const handleClaim = async () => {
    await claimCoupon(1);
    message.success('领取成功');
    fetchCoupons();
  };

  const renderCoupon = (c: CouponVO) => (
    <Card size="small">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <div style={{ fontWeight: 700 }}>
            {c.type === COUPON_TYPE.FULL_REDUCTION ? (
              <span>
                {formatPrice(c.amount)} <span style={{ fontSize: 12, color: '#999' }}>满{c.threshold}可用</span>
              </span>
            ) : (
              <span>
                {formatPercent(c.discount_rate)} <span style={{ fontSize: 12, color: '#999' }}>满{c.threshold}可用</span>
              </span>
            )}
          </div>
          <div style={{ color: '#999', fontSize: 12 }}>{c.name}</div>
          <div style={{ color: '#999', fontSize: 12 }}>
            有效期至 {c.valid_to?.slice(0, 10)}
          </div>
        </div>
        <div style={{ textAlign: 'right' }}>
          <Tag color={c.status === 'UNUSED' ? 'blue' : c.status === 'USED' ? 'default' : 'red'}>
            {c.status === 'UNUSED' ? '未使用' : c.status === 'USED' ? '已使用' : '已过期'}
          </Tag>
        </div>
      </div>
    </Card>
  );

  return (
    <div style={{ maxWidth: 900, margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Typography.Title level={4} style={{ margin: 0 }}>
          我的优惠券
        </Typography.Title>
        <Button type="primary" onClick={handleClaim}>
          领取新人满减券
        </Button>
      </div>
      {coupons.length === 0 ? (
        <EmptyState description="暂无优惠券" />
      ) : (
        <Row gutter={[16, 16]}>
          {coupons.map((c) => (
            <Col key={c.id} xs={24} sm={12} md={8}>
              {renderCoupon(c)}
            </Col>
          ))}
        </Row>
      )}
    </div>
  );
}
