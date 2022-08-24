# coding=utf-8
# *** WARNING: this file was generated by Pulumi SDK Generator. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from . import _utilities

__all__ = [
    'CustomDomainArgs',
    'DatabaseArgs',
    'NetworkingArgs',
]

@pulumi.input_type
class CustomDomainArgs:
    def __init__(__self__, *,
                 domain_name: Optional[pulumi.Input[str]] = None,
                 hosted_zone_name: Optional[pulumi.Input[str]] = None):
        """
        Options for setting a custom domain.
        """
        if domain_name is not None:
            pulumi.set(__self__, "domain_name", domain_name)
        if hosted_zone_name is not None:
            pulumi.set(__self__, "hosted_zone_name", hosted_zone_name)

    @property
    @pulumi.getter(name="domainName")
    def domain_name(self) -> Optional[pulumi.Input[str]]:
        return pulumi.get(self, "domain_name")

    @domain_name.setter
    def domain_name(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "domain_name", value)

    @property
    @pulumi.getter(name="hostedZoneName")
    def hosted_zone_name(self) -> Optional[pulumi.Input[str]]:
        return pulumi.get(self, "hosted_zone_name")

    @hosted_zone_name.setter
    def hosted_zone_name(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "hosted_zone_name", value)


@pulumi.input_type
class DatabaseArgs:
    def __init__(__self__, *,
                 engine_version: Optional[pulumi.Input[str]] = None):
        """
        The options for configuring your database.
        :param pulumi.Input[str] engine_version: The database engine version. Updating this argument results in an outage. See the
               [Aurora MySQL](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Updates.html)
               documentation for your configured engine to determine this value. For example with Aurora MySQL 2,
               a potential value for this argument is 5.7.mysql_aurora.2.03.2. The value can contain a partial version
               where supported by the API.
        """
        if engine_version is None:
            engine_version = '5.7.mysql_aurora.2.08.3'
        if engine_version is not None:
            pulumi.set(__self__, "engine_version", engine_version)

    @property
    @pulumi.getter(name="engineVersion")
    def engine_version(self) -> Optional[pulumi.Input[str]]:
        """
        The database engine version. Updating this argument results in an outage. See the
        [Aurora MySQL](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Updates.html)
        documentation for your configured engine to determine this value. For example with Aurora MySQL 2,
        a potential value for this argument is 5.7.mysql_aurora.2.03.2. The value can contain a partial version
        where supported by the API.
        """
        return pulumi.get(self, "engine_version")

    @engine_version.setter
    def engine_version(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "engine_version", value)


@pulumi.input_type
class NetworkingArgs:
    def __init__(__self__, *,
                 db_subnet_ids: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 ecs_subnet_ids: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 lb_subnet_ids: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None):
        """
        The options for networking.
        :param pulumi.Input[Sequence[pulumi.Input[str]]] db_subnet_ids: The subnets to use for the RDS instance.
        :param pulumi.Input[Sequence[pulumi.Input[str]]] ecs_subnet_ids: The subnets to use for the Fargate task.
        :param pulumi.Input[Sequence[pulumi.Input[str]]] lb_subnet_ids: The subnets to use for the load balancer.
        """
        if db_subnet_ids is not None:
            pulumi.set(__self__, "db_subnet_ids", db_subnet_ids)
        if ecs_subnet_ids is not None:
            pulumi.set(__self__, "ecs_subnet_ids", ecs_subnet_ids)
        if lb_subnet_ids is not None:
            pulumi.set(__self__, "lb_subnet_ids", lb_subnet_ids)

    @property
    @pulumi.getter(name="dbSubnetIds")
    def db_subnet_ids(self) -> Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]:
        """
        The subnets to use for the RDS instance.
        """
        return pulumi.get(self, "db_subnet_ids")

    @db_subnet_ids.setter
    def db_subnet_ids(self, value: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]):
        pulumi.set(self, "db_subnet_ids", value)

    @property
    @pulumi.getter(name="ecsSubnetIds")
    def ecs_subnet_ids(self) -> Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]:
        """
        The subnets to use for the Fargate task.
        """
        return pulumi.get(self, "ecs_subnet_ids")

    @ecs_subnet_ids.setter
    def ecs_subnet_ids(self, value: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]):
        pulumi.set(self, "ecs_subnet_ids", value)

    @property
    @pulumi.getter(name="lbSubnetIds")
    def lb_subnet_ids(self) -> Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]:
        """
        The subnets to use for the load balancer.
        """
        return pulumi.get(self, "lb_subnet_ids")

    @lb_subnet_ids.setter
    def lb_subnet_ids(self, value: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]):
        pulumi.set(self, "lb_subnet_ids", value)


