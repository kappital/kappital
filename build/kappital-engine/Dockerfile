# Copyright 2022 Huawei Cloud Computing Technologies Co., Ltd
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
FROM euleros:latest

WORKDIR /opt/kappital/kappital-engine

RUN groupadd -g 10000 kappital && \
    useradd -u 10000 -g 10000 kappital && \
    echo "Defaults targetpw" >> /etc/sudoers && \
    mkdir -p /opt/kappital/kappital-engine && \
    mkdir -p /opt/kappital/log && \
    touch /opt/kappital/log/kappital-engine.log

COPY bin/kappital-engine /opt/kappital/kappital-engine/

RUN chown -R kappital:kappital /opt/kappital/kappital-engine && \
    chown -R kappital:kappital /opt/kappital/log && \
    chmod 750 /opt/kappital/kappital-engine && \
    chmod 640 /opt/kappital/log/kappital-engine.log && \
    chmod 550 /opt/kappital/kappital-engine/kappital-engine
